package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

//constant buffer size
const BUFFERSIZE = 1024

//single function to rectify all errors
func is_error(err error) bool {
	if err != nil {
		fmt.Println("some error is there", err.Error)
		return true
	}
	return false
}

func main() {
	//server starts listening on port 27001
	//my_address := ""
	if len(os.Args) <= 0 {
		fmt.Println("Enter Socket information")
		os.Exit(0)
		//my_address = os.Args[1]
	}

	listener, err := net.Listen("tcp", os.Args[1])
	is_error(err)

	//to close server automatically
	defer listener.Close()

	fmt.Println("###  Server started! Waiting for connections  ###")
	for {
		//extract connection from listening_port queue
		connection, err := listener.Accept()
		// is_error(err)
		if err != nil {
			fmt.Println("%#v", err)
		}
		//to automatically close client connection
		fmt.Println("Client connected")
		defer connection.Close()

		//to run shell script (schell script takes screenshot)
		_, err = exec.Command("./take_shot/take_shot.sh").Output()
		is_error(err)

		//calling goroutines
		go serve_request(connection)
	}
}

func serve_request(connection net.Conn) {

	// fmt.Println("A client has connected!")

	//opening latest screenshot taken (every recent screenshot is named "snapshot.png")
	//and "./server_copies/" is the path relative to current working directory
	file, err := os.Open("./server_copies/snapshot.png")
	is_error(err)

	//Stat() returns information regarding size,name etc about file opened
	file_info, err := file.Stat()
	is_error(err)

	file_size := fill_buffer(strconv.FormatInt(file_info.Size(), 10), 10)
	file_name := fill_buffer(file_info.Name(), 64)

	// fmt.Println("Sending filename and filesize")
	connection.Write([]byte(file_size))
	connection.Write([]byte(file_name))

	//creating buffer (into which file is read then sent onto "connection")
	sendBuffer := make([]byte, BUFFERSIZE)

	// fmt.Println("Start sending file!")
	for {
		//read content of image file into buffer
		_, err = file.Read(sendBuffer)
		//read  until END_OF_FILE
		if err == io.EOF {
			break
		}
		is_error(err)
		//write content into buffer
		connection.Write(sendBuffer)
	}
	updatelog(connection)

	fmt.Println("File sent  \n\n")

	return
}

//function to fill remaining buffer with addtional characters
func fill_buffer(rec_string string, length int) string {
	for {
		lengtString := len(rec_string)
		if lengtString < length {
			rec_string = rec_string + ":"
			continue
		}
		break
	}
	return rec_string
}

//mainting connection logs at server side
func updatelog(conn net.Conn) {
	// open a file
	f, err := os.OpenFile("./logrecords/record.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// fmt.Println(logrecord[parent])
	if err != nil {
		// fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	log.WithFields(log.Fields{
		"request from: %#v ": conn,
		"at time ":           time.Now().String(),
	}).Info("served")

}

//function to retrive local ip address
// func GetLocalIP() string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return ""
// 	}
// 	for _, address := range addrs {
// 		// check the address type and if it is not a loopback the display it
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}
// 	return ""
// }
