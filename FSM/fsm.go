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
		"Cancel",
		fsm.Events{
			{Name: "Cancel", Src: []string{"Arming", "Armed"}, Dst: "Stop"},
			{Name: "InitArm", Src: []string{"Stop"}, Dst: "Arming"},
			{Name: "Tocheck", Src: []string{"Arming"}, Dst: "system-on-check"},
			{Name: "Checkfailed", Src: []string{"system-on-check"}, Dst: "Arming"},
			{Name: "Arm", Src: []string{"system-on-check"}, Dst: "Armed"},
		},
		fsm.Callbacks{
			"Cancel": func(e *fsm.Event) {
				fmt.Println("Stopping all pod processes: " + e.FSM.Current())
			},
			"InitArm": func(e *fsm.Event) {
				fmt.Println("Setting up the pod: " + e.FSM.Current())
			},
			"Tocheck": func(e *fsm.Event) {
				fmt.Println("Verifying pod functionality: " + e.FSM.Current())
			},
			"Checkfailed": func(e *fsm.Event) {
				fmt.Println("Sensors not working, try to reinitialize " + e.FSM.Current())
			},
			"Arm": func(e *fsm.Event) {
				fmt.Println("Ready for deployment: " + e.FSM.Current())
			},
		},
	)

	fmt.Println(fsm.Current())

	err := fsm.Event("Cancel")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("1 " + fsm.Current())

	err = fsm.Event("InitArm")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("2 " + fsm.Current())

	err = fsm.Event("Tocheck")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("3 " + fsm.Current())

	err = fsm.Event("Arm")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("4 " + fsm.Current())

}
