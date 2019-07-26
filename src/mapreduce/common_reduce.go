package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//
	keys:= make([]string, 0)
	kvs := make(map[string][]string)
	for m := 0; m < nMap; m++{
		reduceFile := reduceName(jobName, m, reduceTask)
		fd, err := os.Open(reduceFile)
		if err != nil{
			log.Fatalf("%s file can't open: %s", reduceFile, err)
		}
		dec := json.NewDecoder(fd)
		var kv KeyValue
		for err := dec.Decode(&kv); err == nil; {
			if _ , ok :=kvs[kv.Key]; !ok{
				keys = append(keys, kv.Key)
			}
			kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
		}
		err = fd.Close()
		if err != nil{
			log.Fatalf("%s file can't close: %s", reduceFile, err)
		}
	}
	sort.Strings(keys)
	output, err := os.Create(outFile)
	if err != nil{
		log.Fatalf("output file can't create: %s", err)
	}
	enc := json.NewEncoder(output)
	for _, key := range keys{
		err = enc.Encode(KeyValue{key, reduceF(key, kvs[key])})
		if err != nil{
			log.Fatalf("file can't encode: %s", err)
		}
	}
	err = output.Close()
	if err != nil{
		log.Fatalf("output file can't close: %s", err)
	}

}
