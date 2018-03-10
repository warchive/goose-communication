package main

import (
	"fmt"
	"github.com/looplab/fsm"
)

//Name: current state name
//Src: previous state that transitions into current state
//Dst: next state that transitions from current state


func main() {
	fsm := fsm.NewFSM(
		"idle"
		fsm.Events{
			{Name: "stop", Src: []string{"idle","STOP","CANCEL"}, Dst: "ARM"},
			{Name: "arming", Src: []string{"ARM","FAILED"}, Dst: "TOCHECK"},
			{Name: "arming", Src: []string{"STOP","FAILED"}, Dst: "STOP"},
			{Name: "system-on-check", Src: []string{"TOCHECK"}, Dst: "SUCCESS"},
			{Name: "system-on-check", Src: []string{"TOCHECK"}, Dst: "FAILED"},
			{Name: "armed", Src: []string{"SUCCESS"}, Dst: "START"},
			{Name: "armed", Src: []string{"SUCCESS"}, Dst: "STOP"},
		},
		fsm.Callbacks{
			"scan": func(e *fsm.Event) {
				fmt.Println("Stopping all pod processes: " + e.FSM.Current())
			},
			"arming": func(e *fsm.Event) {
				fmt.Println("Setting up the pod: " + e.FSM.Current())
			},
			"system-on-check": func(e *fsm.Event) {
				fmt.Println("Verifying pod functionality: " + e.FSM.Current())
			},
			"armed": func(e *fsm.Event) {
				fmt.Println("Ready: " + e.FSM.Current())
			},
		},
	)

	fmt.Println(fsm.Current())

	err := fsm.Event("scan")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("1:" + fsm.Current())

	err = fsm.Event("working")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("2:" + fsm.Current())

	err = fsm.Event("situation")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("3:" + fsm.Current())

	err = fsm.Event("finish")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("4:" + fsm.Current())

}
