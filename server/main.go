package main

import (
	//"net"
	//"context"
	//"log"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	//"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/dynamicpb"
	//"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

	//"google.golang.org/protobuf/reflect/protoreflect"

	//pref "google.golang.org/protobuf/reflect/protoreflect"
	//preg "google.golang.org/protobuf/reflect/protoregistry"
	apiPB "github.com/michaelhenkel/sub/api/proto"
	serverPB "github.com/michaelhenkel/sub/server/proto"
)

type server struct {
	serverPB.UnimplementedServerServer
}

//go:generate protoc -I../policy/proto -I$GOPATH/src/github.com/gogo/protobuf/gogoproto --proto_path=proto --gogo_out=plugins=grpc:proto --gogo_opt=paths=source_relative --include_source_info --include_imports --descriptor_set_out proto/genbyte/desc.protoset proto/server.proto

func newServer() *server {
	s := &server{}
	return s
}

type myFileDescriptor struct{}

func main() {
	fds := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(DSCByte, fds)
	if err != nil {
		fmt.Println(err)
	}

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	client := apiPB.NewApiClient(conn)

	for _, d := range fds.GetFile() {
		if *d.Name == "gogo.proto" {
			msg := &apiPB.Message{
				Fd: d,
			}
			_, err = client.RegisterType(context.Background(), msg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, d := range fds.GetFile() {
		if *d.Name == "server.proto" {
			msg := &apiPB.Message{
				Fd: d,
			}
			_, err = client.RegisterType(context.Background(), msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

/*
	resyncPeriod := time.Second * 60
	var myObj runtime.Object
	myObj = &MyObject{}
	cc := &cache.Config{
		ObjectType: myObj,
		ListerWatcher: &cache.ListWatch{
			ListFunc:  myListFunc,
			WatchFunc: myWatchFunc,
		},
		Queue:            &MyQueue{},
		Process:          myProcessFunc,
		FullResyncPeriod: resyncPeriod,
		RetryOnError:     true,
	}
	cache.SharedInformer
	var sII cache.SharedIndexInformer
	myStore := &MyStore{}
	mySharedIndexInformer := &MySharedIndexInformer{}
	sII = mySharedIndexInformer
	cacheController := cache.New(cc)
	stopChan := make(<-chan struct{})
	cacheController.Run(stopChan)
	<-stopChan

}

func myProcessFunc(obj interface{}) error {
	return nil
}

func myWatchFunc(options metav1.ListOptions) (watch.Interface, error) {
	myWatch := &MyWatch{}
	return myWatch, nil

}

func myListFunc(options metav1.ListOptions) (runtime.Object, error) {
	myObj := &MyObject{}
	return myObj, nil
}

type MyStore struct {
}

func (ms *MyStore) Add(obj interface{}) error {
	return nil
}

func (ms *MyStore) Update(obj interface{}) error {
	return nil
}

func (ms *MyStore) Delete(obj interface{}) error {
	return nil
}

func (ms *MyStore) List() []interface{} {
	var intfList []interface{}
	return intfList
}

func (ms *MyStore) ListKeys() []string {
	var keyList []string
	return keyList
}

func (ms *MyStore) GetByKey(key string) (item interface{}, exists bool, err error) {
	return item, exists, err
}

func (ms *MyStore) Replace([]interface{}, string) error {
	return nil
}

func (ms *MyStore) Resync() {
}

type MyIndexer struct {
	cache.Store
}

func (mi *MyIndexer) Index(indexName string, obj interface{}) ([]interface{}, error) {
	var intfList []interface{}
	return intfList, nil
}

func (mi *MyIndexer) IndexKeys(indexName, indexKey string) ([]string, error) {
	var indexKeyList []string
	return indexKeyList, nil
}

func (mi *MyIndexer) ListIndexFuncValues(indexName string) []string {
	var indexKeyList []string
	return indexKeyList
}

func (mi *MyIndexer) ByIndex(indexName, indexKey string) ([]interface{}, error) {
	var intfList []interface{}
	return intfList, nil
}

func (mi *MyIndexer) GetIndexers() cache.Indexers {
	var cacheIndexers cache.Indexers
	return cacheIndexers
}

func (mi *MyIndexer) AddIndexers(newIndexers cache.Indexers) error {
	return nil
}

type MySharedIndexInformer struct {
}

func (msi *MySharedIndexInformer) AddIndexers(indexers cache.Indexers) error {
	return nil
}

func (msi *MySharedIndexInformer) GetIndexer() cache.Indexer {
	return nil
}

type MyQueue struct {
	cache.Store
}

func (mq *MyQueue) Pop(cache.PopProcessFunc) (interface{}, error) {
	var intf interface{}
	return intf, nil
}

func (mq *MyQueue) AddIfNotPresent(interface{}) error {
	return nil
}

func (mq *MyQueue) HasSynced() bool {
	return true
}

func (mq *MyQueue) Close() {
}

type MyWatch struct {
}

func (mw *MyWatch) Stop() {

}

func (mw *MyWatch) ResultChan() <-chan watch.Event {
	eventChan := make(<-chan watch.Event)
	return eventChan

}

type MyObject struct {
}

type MyObjectKind struct {
}

func (m *MyObject) GetObjectKind() schema.ObjectKind {
	mk := &MyObjectKind{}
	return mk
}

func (m *MyObject) DeepCopyObject() runtime.Object {
	return m
}

func (mk *MyObjectKind) SetGroupVersionKind(kind schema.GroupVersionKind) {

}

func (mk *MyObjectKind) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{}

}
*/
