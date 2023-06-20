package main

import (
	"bufio"
	"fmt"
	"os"
)

type EventType byte

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type TransactionLogger interface {
	WritePut(key, value string)
	WriteDelete(key string)
	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
}

type FileTransactionLogger struct {
	Events       chan<- Event // 송신 이벤트에 대한 쓰기 전용 채널
	Errors       <-chan error // 수신 이벤트에 대한 읽기 전용 채널
	LastSequence uint64       // 마지막으로 사용한 이벤트 일련 번호
	File         *os.File     // 트랜잭션 로그의 위치
}

func (ftl *FileTransactionLogger) WritePut(key, value string) {
	ftl.Events <- Event{
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (ftl *FileTransactionLogger) WriteDelete(key string) {
	ftl.Events <- Event{EventType: EventDelete, Key: key}
}

func (ftl *FileTransactionLogger) Err() <-chan error {
	return ftl.Errors
}

func (ftl *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	ftl.Events = events

	errors := make(chan error, 1)
	ftl.Errors = errors

	go func() {
		for e := range events {
			ftl.LastSequence++
			_, err := fmt.Fprintf(ftl.File, "%d\t%d\t%s\t%s\n", ftl.LastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (ftl *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(ftl.File)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if ftl.LastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			ftl.LastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func NewFileTransactionLogger(fileName string) (TransactionLogger, error) {
	/*
		os.O_RDWR		 파일을 읽고 쓸 수 있는 모드로 연다.
		os.O_APPEND		 파일에 대한 쓰기 요청은 파일 뒷부분에 추가되며 기존 내용을 덮어쓰지 않는다.
		os.O_CREATE		 파일이 존재하지 않는 경우 생성하고 연다.
	*/
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{File: file}, nil
}
