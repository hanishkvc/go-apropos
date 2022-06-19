// Test some aspects of Go interfaces
// HanishKVC, 2022

package main

import (
	"fmt"
	"regexp"
)

type IFMyInt int
type IFBaseStruct struct {
	f1 int
}
type IFMyStruct IFBaseStruct
type IFMyRE regexp.Regexp

type IFMagic interface {
	whome() string
	other() string
}

func (p *IFMyInt) whome() string {
	return "IFMyInt"
}

func (v IFMyInt) other() string {
	return "Other-IFMyInt"
}

func (p IFMyStruct) whome() string {
	return "IFMyStruct"
}

func (v IFMyStruct) other() string {
	return "Other-IFMyStruct"
}

func (p *IFMyRE) whome() string {
	return "IFMyRE"
}

/* Cant create a method for a pointer type

type IFMyIntPtr *int

func (v IFMyIntPtr) cani() string {
	return "IFMyIntPtr"
}

*/

func test_if() {
	var ifMyInt IFMyInt
	var baseInt int
	var ifMagicInt IFMagic
	baseIntPtr := &baseInt

	ifMyInt = IFMyInt(baseInt)
	ifMyInt = IFMyInt(*baseIntPtr)

	/* 50-50
	ifMagicInt = IFMyInt(baseInt)
	ifMagicInt = IFMyInt(*baseIntPtr)
	*/

	ifMagicInt = &ifMyInt
	ifMagicInt = IFMagic(&ifMyInt)

	/* Never ok
	ifMagicInt = baseInt
	ifMagicInt = &baseInt
	ifMagicInt = IFMagic(baseInt)
	ifMagicInt = IFMagic(&baseInt)
	ifMagicInt = IFMagic(baseIntPtr)
	ifMagicInt = IFMagic(*baseIntPtr)
	*/

	fmt.Printf("%v:INFO:T IF: ifMyInt:%v\n", PRG_TAG, ifMyInt)
	fmt.Printf("%v:INFO:T IF: ifMagicInt:%v\n", PRG_TAG, ifMagicInt)

	var baseStruct IFBaseStruct
	var ifMyStruct IFMyStruct
	baseStructPtr := &baseStruct
	ifMyStruct = IFMyStruct(baseStruct)
	ifMyStruct = IFMyStruct(*baseStructPtr)
	var ifMagicStruct = IFMagic(ifMyStruct)
	ifMagicStruct = IFMagic((*IFMyStruct)(&baseStruct))
	fmt.Printf("%v:INFO:T IF: ifMyStruct:%v\n", PRG_TAG, ifMyStruct)
	fmt.Printf("%v:INFO:T IF: ifMagicInt:%v\n", PRG_TAG, ifMagicStruct)
}
