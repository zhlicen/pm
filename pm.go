package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"pm/ssh"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type MonitorRecord struct {
	ProcessName string
	Time        string
	CpuUsage    string
	VmSize      string
	VmRss       string
}

func ParseResult(result string) (record MonitorRecord, err error) {

	r, _ := regexp.Compile(`([0-9\.]+)\s+(\d+)\s+(\d+)+\w*`)
	if r.MatchString(result) {
		matchResult := r.FindAllStringSubmatch(result, -1)
		if matchResult == nil {
			err = errors.New("Invalid result: " + result)
			return
		}
		var record MonitorRecord
		record.CpuUsage = matchResult[0][1]
		record.VmSize = matchResult[0][2]
		record.VmRss = matchResult[0][3]
		log.Printf("Parse result: cpu=%s, vmsize=%s, vmrss=%s\n",
			record.CpuUsage, record.VmSize, record.VmRss)
		return record, nil
	}
	err = errors.New("Invalid result: " + result)
	return
}

// WriteRecordToInfluxDB Write record to influx db
func WriteRecordToInfluxDB(address, db, measurement string, record MonitorRecord) {
	url := fmt.Sprintf("%swrite?db=%s", address, db)
	format := "%s,process=%s cpu=%s,vmsize=%s,vmrss=%s %s"
	body := fmt.Sprintf(format, measurement,
		record.ProcessName, record.CpuUsage, record.VmSize, record.VmRss, record.Time)
	log.Printf("POST " + url + "\n" + body)
	rsp, err := http.Post(url, "application/octet-stream", strings.NewReader(body))
	if err != nil {
		log.Println(err)
	}
	log.Println("returns:" + rsp.Status)
}

func main() {
	// Read configuration file
	viper.SetConfigName("pm")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	sshIP := viper.GetString("ssh.ip")
	sshPort := viper.GetString("ssh.port")
	sshUserName := viper.GetString("ssh.username")
	sshPassword := viper.GetString("ssh.password")
	monitorProcesses := viper.GetStringSlice("global.processes")
	loopPerMinutes := viper.GetInt("global.loop_every_minutes")
	dbAddress := viper.GetString("influxdb.address")
	dbDb := viper.GetString("influxdb.db")
	dbMesurement := viper.GetString("influxdb.measurement")
	delay := time.Duration(loopPerMinutes) * time.Minute
	for {
		sshClient := new(utility.SSHClient)
		err = sshClient.Connect(sshIP, sshPort, sshUserName, sshPassword)
		if err != nil {
			log.Println(err)
			time.Sleep(delay)
			continue
		}
		for _, process := range monitorProcesses {
			log.Println("Reading process information for: ", process)
			command := fmt.Sprintf(`ps -axo cmd,pcpu,rss,vsz | grep %s | grep -v "grep"`, process)
			result := sshClient.Run(command)
			if result != "" {
				fmt.Println(result)
				record, err := ParseResult(result)
				if err == nil {
					record.ProcessName = process
					timeCommand := "date +%s%N"
					record.Time = sshClient.Run(timeCommand)
					WriteRecordToInfluxDB(dbAddress, dbDb, dbMesurement, record)
				} else {
					log.Println(err)
				}
			} else {
				log.Printf("Command: [%s] \nreturns nothing\n", command)
			}
		}
		sshClient.Disconnect()
		time.Sleep(delay)
	}
}
