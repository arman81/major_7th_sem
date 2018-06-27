package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andlabs/ui"
)

//constant for buffer_size (must be of same size as that of server)
const BUFFERSIZE = 1024

func main() {
	err := ui.Main(func() {
		address := ui.NewEntry()
		button := ui.NewButton("Get Shot")
		greeting := ui.NewLabel("")
		port := ui.NewEntry()
		starting_address := ui.NewEntry()
		cport := ui.NewEntry()
		dropdown := ui.NewCombobox()
		getaddress := ui.NewButton("Available Monitors")

		box := ui.NewVerticalBox()

		box.Append(ui.NewLabel("Dropdown Addressed / Manually Insert"), false)
		//code to update values in dropdown_menu

		box.Append(starting_address, false)
		box.Append(cport, false)
		box.Append(dropdown, false)
		box.Append(getaddress, false)

		box.Append(ui.NewLabel(""), false)
		box.Append(ui.NewLabel(""), false)

		box.Append(ui.NewLabel("Enter IP address:"), false)
		box.Append(address, false)
		box.Append(ui.NewLabel("Enter port :"), false)
		box.Append(port, false)
		box.Append(button, false)
		box.Append(greeting, false)

		window := ui.NewWindow("Get Shot", 600, 500, false)
		window.SetChild(box)

		//map to remove duplicate from dropdown_menu
		available_addresses := make(map[string]int)
		//BUTTON for updating dropdown_menu
		getaddress.OnClicked(func(*ui.Button) {
			//clicked to insert values in dropdown_menu
			// _, err := exec.Command("./check_con/is_open.sh").Output()
			// start_address := "127.0.0"
			// eport := "23001"

			f, err := os.OpenFile("./check_con/available_con.txt", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				greeting.SetText("Error in checking available ports")
			}

			defer f.Close()

			for i := 0; i < 255; i++ {
				temp_address := string(starting_address.Text()) + "." + strconv.Itoa(i) + ":" + string(cport.Text())
				tcpAddr, err := net.ResolveTCPAddr("tcp4", temp_address)
				// is_error(err)
				if err != nil {
					continue
				}

				_, err = net.DialTCP("tcp", nil, tcpAddr)
				// is_error(err)
				if err != nil {
					continue
				}
				if _, err = f.WriteString(temp_address + "\n"); err != nil {
					greeting.SetText("Error while updating list")
				}
			}

			if err != nil {
				greeting.SetText("No address is available right now due to some technical error")
			} else {

				//reading from file to update list
				file, err := os.Open("./check_con/available_con.txt")
				if err != nil {
					greeting.SetText("Error while getting available list (file_opeing error)")
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					// greeting.SetText(string(scanner.Text()))
					//check if address is already available in dropdown_menu
					if _, ok := available_addresses[string(scanner.Text())]; !ok {
						available_addresses[string(scanner.Text())] = len(available_addresses)
						dropdown.Append(string(scanner.Text()))
					}
				}

				if err := scanner.Err(); err != nil {
					greeting.SetText("Error while getting available list")
				} else {
					greeting.SetText("Successfully loaded list")
				}
			}
		})

		//handler to deal with selected value of dropdown_menu
		dropdown.OnSelected(func(*ui.Combobox) {
			index := dropdown.Selected()
			if index == -1 {
				greeting.SetText("Select some value from dropdown!")
			} else {
				//function to split host and port from given string "127.0.0.1:23001"
				for key, val := range available_addresses {
					// fmt.Println(strconv.Itoa(index) + "  " + strconv.Itoa(val))
					if val == index {
						host, prt, err := net.SplitHostPort(key)
						if err != nil {
							greeting.SetText("Error while loading values from dropdown to textfields")
						} else {
							address.SetText(host)
							port.SetText(prt)
							greeting.SetText("Address selected: " + strconv.Itoa(index))
						}
						break
					}
				}

			}
		})

		//button ONCLICK function (GET SHOT BUTTON)
		button.OnClicked(func(*ui.Button) {

			//sending ip address to function by converting it into string
			if verify(address.Text(), port.Text()) {
				done := send_request(string(address.Text() + ":" + port.Text()))

				if done == true {
					greeting.SetText("Image received!")
				} else {
					greeting.SetText("Connection error!")
				}

			} else {
				greeting.SetText("Please Provide Valid Inputs!! ")
				address.SetText("")
				port.SetText("")
			}
		})

		//on closing application windows
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		//function to start (open ) windows
		window.Show()
	})
	if err != nil {
		panic(err)
	}

}

//function to send request and handling receiving images
func send_request(address string) bool {

	//flag to return true if image is received (return flag to main())
	//done := false

	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	// is_error(err)
	if err != nil {
		return false
	}

	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	// is_error(err)
	if err != nil {
		return false
	}

	// fmt.Println("successfully dialed connection")

	//message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Message from server: " + message)

	//buffer to read file_name
	buffer_file_name := make([]byte, 64)

	//buffer to read file_size
	buffer_file_size := make([]byte, 10)

	connection.Read(buffer_file_size)
	//converting received file_size into 64bitInt (after trimming addtional ":"s)
	file_size, _ := strconv.ParseInt(strings.Trim(string(buffer_file_size), ":"), 10, 64)

	connection.Read(buffer_file_name)
	//trimming additional ":"s from file_name
	file_name := strings.Trim(string(buffer_file_name), ":")

	// fmt.Println("Received filesize =   ", file_size, "  and file_name =   ", file_name)

	//save received image as SOCKET+TIMESTAMP
	// fmt.Println("creating newFile object")
	newFile, err := os.Create("./received/" + address + "_" + string(file_name) + time.Now().String())

	is_error(err)

	defer newFile.Close()

	//variable to store TOTAL_RECEIVED_BYTES
	var receivedBytes int64

	for {

		//if remaining bytes are less then BUFFERSIZE(last byte)
		if (file_size - receivedBytes) < BUFFERSIZE {
			//copy content of buffer to newFile
			io.CopyN(newFile, connection, (file_size - receivedBytes))
			//clearing buffer for further operations
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-file_size))
			break
		}
		//copying content to newFile
		io.CopyN(newFile, connection, BUFFERSIZE)
		//increase counter of TOTAL_RECEIVED_BYTES
		receivedBytes += BUFFERSIZE
	}

	// fmt.Println("Received file completely!")
	return true
}

//commmon function to handle errors
func is_error(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s,%#v", err.Error(), err)
		os.Exit(1)
	}
}

//function to validate ip address and port number provided
func verify(address string, port string) bool {
	addr := net.ParseIP(address)
	if addr == nil {
		return false
	}

	p, err := strconv.Atoi(port)

	if p < 0 || p > 65536 || err != nil {
		return false
	}

	return true
}
