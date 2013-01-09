/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */

package TAppEncoder

import (
    "io"
    "bufio"
    "fmt"
    //"errors"
    "os"
    "log"
    "strings"
	"strconv"
)

type ParseFailure struct {
  Arg   string;
  Val   string;
};

type Option interface{
    GetName() string;
	Parse(arg string);
	SetDefault();
}

/** Type specific option storage */
//template<typename T>
type OptionBase struct{
	opt_name        	string;
  	opt_desc        	string; 
    opt_storage     	interface{};
    opt_default_value 	interface{};
}

func (this *OptionBase) GetName() string{
    return this.opt_name;
}

func (this *OptionBase) SetDefault() {
    this.opt_storage = this.opt_default_value;
}

type OptionString struct{
	OptionBase
}

func NewOptionString(name string, storage interface{}, default_value interface{}, desc string) *OptionString{
    return &OptionString{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionString) Parse(arg string) {
   this.opt_storage = arg;
}


type OptionInt struct{
	OptionBase
}

func NewOptionInt(name string, storage interface{}, default_value interface{}, desc string) *OptionInt{
    return &OptionInt{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionInt) Parse(arg string) error {
	argint, err := strconv.Atoi(arg);
	if err!=nil{
		return err
	}
	
	this.opt_storage = argint;
	
	return nil;
}
   
type Options struct {
  opt_map 	map[string]Option; //std::list<Names*>
};

func NewOptions() *Options {
	return &Options{opt_map:make(map[string]Option)};
}

func (opts *Options) AddOption (opt Option){
	  opts.opt_map[opt.GetName()] = opt;
}

func (opts *Options) DoHelp    (columns uint){
}


func (opts *Options) StorePair( name string, value string) bool{
      //Options::NamesMap::iterator opt_it;
      opt, found := opts.opt_map[name];

      if !found {
        // not found
        fmt.Printf("Unknown option: '%s' (value:`%s')\n", name, value);
        return false;
      }

      opt.Parse(value);
      
      return true;
}

func (opts *Options) ScanLine  ( line string){
      /* strip any leading whitespace */
      line = strings.TrimSpace(line);
      if line == "" {
        /* blank line */
        return;
      }
      if line[0:1] == "#" {
        /* comment line */
        return;
      }
      commentPos := strings.Index(line, "#");
      if commentPos!=-1 {
      	line = line[0:commentPos];
      }
      commaPos := strings.Index(line, ":");
      if commaPos== -1 {
      	// error: badly formatted line
        return;
      }
      name  := line[0:commaPos];
      value := line[commaPos+1:];

      /* store the value in option */
      opts.StorePair(name, value);
}
func (opts *Options) ScanFile  ( in io.Reader) (err error){
	var line string
	reader := bufio.NewReader(in)
    eof := false

    line, err = reader.ReadString('\n')
    if err == io.EOF {
        err = nil
        eof = true
    } else if err != nil {
        return err
    }

    for !eof {
    	opts.ScanLine(line);
    	
    	line, err = reader.ReadString('\n')
	    if err == io.EOF {
	        err = nil
	        eof = true
	    } else if err != nil {
	        return err
	    }
    }
    
    return nil;
}

func (opts *Options) ParseConfigFile( filename string){
	cfgstream, err := os.Open(filename);
	if err!=nil {
		log.Fatal(err)
	}
	defer cfgstream.Close()
	
	opts.ScanFile(cfgstream);
}

