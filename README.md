# go-uds
Go implementation of UDS (ISO-14229) standard.

This project is an implementation of the Unified Diagnostic Services (UDS) protocol defined by ISO-14229 written in golang.

It is a port of python version: https://github.com/pylessard/python-udsoncan

UDS runs on top of the isotp protocol:

also port of Python version https://github.com/pylessard/python-can-isotp

does not require https://github.com/hartkopp/can-isotp

but reading canbus data from real device requires native c code for darwin, windows and linux.

you define stack_rxfn and stack_txfn functions that will call the native c code for real data.



