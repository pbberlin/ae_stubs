package blobstore_image_resize

import (
	"net/http"
	"fmt"
	"bytes"
	"math"
	
	//"github.com/pbberlin/tools/util"
	//"github.com/pbberlin/tools/u_err"

)


var scale_x_vm,scale_y_vm           map[float64][]float64  // vector - map

func init() {

	scale_y_vm = map[float64][]float64{
		 7.5 : nil,
		 5.0 : nil,
		 2.5 : nil,
		 1.0 : nil,
	}
	fillScale(scale_y_vm,5)


	scale_x_vm = map[float64][]float64{
		 9.0 : nil,
		 6.0 : nil,
		 2.4 : nil,
		 1.2 : nil,
	}
	fillScale(scale_x_vm,12)
	
	http.HandleFunc("/blob/scale"  , blobScale)
	
	
}

func getScale( f_max float64, scales map[float64][]float64 )( keyScale float64, potency int, msg string){

	b1  := new(bytes.Buffer) 
	b1.WriteString( "<hr>\n" )
	
	ff_max := float64(f_max)
	b1.WriteString( fmt.Sprintf("searching maxval for %#v<br>\n", ff_max))
	
	tmp := ff_max
	if tmp < 0 { tmp = tmp * -1}
	if tmp >= 1 {
		for {
			tmp = tmp / 10
			if tmp < 1 {break}
			potency++
		}		
	} else {
		potency = -1
		for {
			tmp = tmp * 10
			if tmp > 1 {break}
			potency--
		}			
	}
	mantisse := ff_max / ( math.Pow10(potency))

	b1.WriteString( fmt.Sprintf("mantisse <b>%6.2f</b> - potency  %#v<br>\n", mantisse, potency ) )


	smallest_dist  := 10.0
	for max_scale_val, _ := range scales {
		lp_dist := max_scale_val - mantisse
		if lp_dist >= 0 && lp_dist <= smallest_dist    {
			keyScale = max_scale_val
			smallest_dist  = lp_dist
		}		
		b1.WriteString( fmt.Sprintf("  &nbsp;  &nbsp; cealinged by %4.2f  --- dist %4.2f - new min %4.2f<br>\n", max_scale_val,lp_dist,smallest_dist) )
	}
	b1.WriteString( fmt.Sprintf("found scale <b>%#v</b> in between<br>\n", keyScale) )


	// skip over
	smallest_scale := 10.0
	largest_scale  :=  0.0
	for max_scale_val, _ := range scales {
		if max_scale_val < smallest_scale {  smallest_scale = max_scale_val}
		if max_scale_val > largest_scale  {  largest_scale  = max_scale_val}		
	}
	if 10* smallest_scale > mantisse  &&  mantisse > largest_scale {
		keyScale = smallest_scale
		b1.WriteString( fmt.Sprintf("<br>found scale <b>%#v</b> - in loop over\n", keyScale) )
		b1.WriteString( fmt.Sprintf(" &nbsp;  &nbsp;  &nbsp; -- smallest scale %#v - largest scale  %#v -- <br>\n", smallest_scale, largest_scale) )
	}

	
	msg = b1.String()
	return 
	
}

func blobScale(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "text/html")

	//fmt.Fprint(w,  util.PrintMap(util.Map_example_right))	
	
	_,_,msg := getScale(110.94, scale_x_vm)
	fmt.Fprintf(w, msg )

	_,_,msg = getScale(0.0094, scale_x_vm)
	fmt.Fprintf(w, msg )

	_,_,msg = getScale(5555.0094, scale_x_vm)
	fmt.Fprintf(w, msg )

	_,_,msg = getScale(9, scale_x_vm)
	fmt.Fprintf(w, msg )

	_,_,msg = getScale(120, scale_x_vm)
	fmt.Fprintf(w, msg )

	fmt.Fprintf(w, printScale(scale_y_vm)  )
	fmt.Fprintf(w, printScale(scale_x_vm)  )




}

func printScale(s map[float64][]float64 ) string{

	b1  := new(bytes.Buffer) 
	b1.WriteString( "<hr>\n" )
	for max_val,vs := range s {
		b1.WriteString( fmt.Sprint("<b>",max_val,"</b><br>\n" ) )
		for i, val :=  range vs{

			quot  := fmt.Sprintf("%-4.2f", val)
			if len(quot)>0  && quot[len(quot)-1:] == "0" {
				quot = quot[:len(quot)-1]
			}
			b1.WriteString( fmt.Sprintf("<pre style='margin:0'> %-6d   %s</pre>\n",i,quot) ) 

		}
		b1.WriteString( "<br>\n" )
	}
	return b1.String()
}

func fillScale( s map[float64][]float64 , num_ticks int){

	for max_val,vs := range s {
		vs  = make( []float64,  num_ticks+1)
		for i:=0;i<=num_ticks;i++ {
			ftick := max_val/float64(num_ticks)
			ftick_val := float64(i) * ftick 
			vs[i] = ftick_val 
		}
		s[max_val] = vs  // unclear why this is neccessary
	}
}

