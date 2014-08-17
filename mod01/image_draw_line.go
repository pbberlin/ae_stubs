package main

import (
	"appengine"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	// _ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"image/color"	
	

	"math"
	"bytes"
	"net/http"
	"io"
	"os"
	"log"

	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/u_err"
	
)



// NOT according to http://en.wikipedia.org/wiki/Alpha_compositing
//   but by my own trial and error
func pixOverPix(ic_old,ic_new uint8, al_new float64)(c_res uint8) {

	al_old := float64(1); _=al_old
	c_old  := float64(ic_old)
	c_new  := float64(ic_new)

	algo1 := c_old*(1-al_new)   +   c_new*al_new
	c_res =  uint8( util.Min( util.Round(algo1),255) )
	//log.Printf("\t\t %3.1f +  %3.1f  = %3.1f", c_old*(1-al_new),c_new*al_new, algo1)

	return 
}



func funcSetPixler(col color.RGBA, img *image.RGBA )( func( addr int,dist float64) ){

	// 	4*400*300
	r  := img.Rect
	p0 := r.Min
	p1 := r.Max
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	maxPix := dx * dy * 4
	log.Printf("\tfuncSetPixler  BxH: %vx%v  Size:%v (%v)",dx,dy,maxPix,len(img.Pix) )	
	
	return func(addr int,dist float64){
		
		//log.Printf("\t%v<%v",addr,maxPix )	
		if addr > (maxPix-4)  ||  addr < 0 {
			log.Printf("\t%v<%v !  OVERFLOW! ",addr,maxPix )	
			return
		}

		// dist ranges from 0 to 1.5
		if dist < 0.0 { 
			dist = 0
		}
		
		sharpness := 0.9  // < 1 => more blurred ; otherwise more pixely
		ba :=  math.Pow( 1 - (dist * 2/3), sharpness )

		//log.Printf("\tbef: %3d %3d %3d",img.Pix[addr+0],img.Pix[addr+1],img.Pix[addr+2])
		//log.Printf("\tcol: %3d %3d %3d | %1.3f => %1.3f",col.R,col.G,col.B,dist,ba)

		img.Pix[addr+0] = pixOverPix(img.Pix[addr+0],col.R,  ba) 
		img.Pix[addr+1] = pixOverPix(img.Pix[addr+1],col.G,  ba) 
		img.Pix[addr+2] = pixOverPix(img.Pix[addr+2],col.B,  ba) 

		//log.Printf("\taft: %3d %3d %3d\n\n",img.Pix[addr+0],img.Pix[addr+1],img.Pix[addr+2])

	}
	
}







// https://courses.engr.illinois.edu/ece390/archive/archive-f2000/mp/mp4/anti.html
func funcDrawLiner(lCol color.RGBA, img *image.RGBA )( func( P_next image.Point, lCol color.RGBA, img *image.RGBA )  ){

	var P_last image.Point = image.Point{-1111,-1111}

	r  := img.Rect
	p0 := r.Min
	p1 := r.Max
	imgWidth := p1.X - p0.X


	return func (P_next image.Point, lCol color.RGBA, img *image.RGBA ){

		var P0, P1 image.Point

		if P_last.X == -1111  &&  P_last.Y == -1111{
			P_last = P_next
			return	
		} else {
			P0 = P_last	
			P1 = P_next
			P_last = P_next	
		}
		
		
		log.Printf("draw_line_start---------------------------------")
	
		x0, y0 := P0.X, P0.Y
		x1, y1 := P1.X, P1.Y
	
		
		bpp := 4  // bytes per pixel
	
		addr := (y0*imgWidth+x0)*bpp
		dx   := x1-x0
		dy   := y1-y0
	
	
		var du, dv,u ,v int
		var uincr int = bpp
		var vincr int = imgWidth*bpp
	
		
		// switching to (u,v) to combine all eight octants
		if  util.Abs(dx) > util.Abs(dy) {
			du = util.Abs(dx)
			dv = util.Abs(dy)
			u = x1
			v = y1
			uincr = bpp
			vincr = imgWidth*bpp
			if dx < 0 {uincr = -uincr}
			if dy < 0 {vincr = -vincr}
		} else {
			du = util.Abs(dy)
			dv = util.Abs(dx)
			u = y1
			v = x1
			uincr = imgWidth*bpp
			vincr = bpp
			if dy < 0 {uincr = -uincr}
			if dx < 0 {vincr = -vincr}
		}
		log.Printf("draw_line\tu %v - v %v - du %v - dv %v - uinc %v - vinc %v ", u, v, du, dv, uincr, vincr)
		
		// uend	  :=  u + 2 * du
		// d	     := (2 * dv) - du		// Initial value as in Bresenham's 
		// incrS   :=  2 *  dv				// Δd for straight increments 
		// incrD   :=  2 * (dv - du)	   // Δd for diagonal increments 
		// twovdu  :=  0						// Numerator of distance starts at 0 
	
	
		// I have NO idea why - but unless I use -1- 
		//   instead of the orginal -2- as factor,
		//   all lines are drawn DOUBLE the intended size
		//   THIS is how it works for me:
		uend	  :=  u + 1 * du
		d	     := (1 * dv) - du		// Initial value as in Bresenham's 
		incrS   :=  1 *  dv				// Δd for straight increments 
		incrD   :=  1 * (dv - du)	   // Δd for diagonal increments 
		twovdu  :=  0						// Numerator of distance starts at 0 
	
						
	
		log.Printf("draw_line\tuend %v - d %v - incrS %v - incrD %v - twovdu %v", uend, d, incrS, incrD, twovdu)
	
	
		tmp     := float64(du*du + dv*dv)
		invD	  := 1.0 / (2.0*math.Sqrt( tmp ))   /* Precomputed inverse denominator */
		invD2du := 2.0 * (  float64(du)*invD)	   /* Precomputed constant */
	
		log.Printf("draw_line\tinvD %v - invD2du %v", invD, invD2du)
	
		cntr := -1
		
		setPixClosure := funcSetPixler(lCol,img)
		
		for{
			cntr++
			//log.Printf("==lp%v ", cntr )
	
			// Ensure that addr is valid
			ftwovdu:= float64(twovdu)
			setPixClosure(addr - vincr, invD2du + ftwovdu*invD)
			setPixClosure(addr		  ,	        ftwovdu*invD)
			setPixClosure(addr + vincr, invD2du - ftwovdu*invD)
		
	
			if (d < 0){
				/* choose straight (u direction) */
				twovdu = d + du
				d = d + incrS
			} 	else 	{
				/* choose diagonal (u+v direction) */
				twovdu = d - du
				d = d + incrD
				v = v+1
				addr = addr + vincr
			}
			u = u+1
			addr = addr+uincr
			
			if u > uend {break}
		} 
	
		log.Printf("draw_line_end---------------------------------")
	
	}
}	





func imageAnalyze(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	
	// This example demonstrates decoding a JPEG image 
	// and examining its pixels.

	// Decode the JPEG data. 
	
	// If reading from file, create a reader with
	// reader, err := os.Open("testdata/video-001.q50.420.jpeg")
	// if err != nil {  c.Errorf(err)  }
	// defer reader.Close()
	

	img,whichformat := conv.Base64_str_to_img(conv.Img_jpeg_base64)	
	c.Infof( "retrieved img from base64: format %v - type %T\n" , whichformat, img )

	bounds := img.Bounds()

	// Calculate a 16-bin histogram for m's red, green, blue and alpha components.
	//
	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	var histogram [16][4]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 12 reduces this to the range [0, 15].
			histogram[r>>12][0]++
			histogram[g>>12][1]++
			histogram[b>>12][2]++
			histogram[a>>12][3]++
		}
	}

	// Print the results.
	b1 := new(bytes.Buffer)

	s1 := fmt.Sprintf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	b1.WriteString( s1 )

	for i, x := range histogram {
		s1 := fmt.Sprintf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
		b1.WriteString( s1 )
	}

	
	w.Header().Set("Content-Type", "text/plain")
	w.Write( b1.Bytes() )
}



func base64_from_file(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	p := r.FormValue("p")
	if p == "" { p = "static/chartbg_400x960__480x1040__12x10.png" }	

	f, err := os.Open(p)

	util_err.Err_http(w,r,err)

	defer f.Close()
	
	img, whichformat, err := image.Decode(f)
	util_err.Err_http(w,r,err)
	c.Infof( "format %v - type %T\n" , whichformat, img )
	imgRGBA,ok := img.(*image.RGBA)
	util_err.Err_http(w,r,ok)

	str_b64 := conv.Rgba_img_to_base64_str( imgRGBA )

	w.Header().Set("Content-Type", "image/plain")
	w.Header().Set("Content-Disposition", "inline;filename=img_as_base64.txt")	
	//w.Header().Set("Content-Disposition", "attachment;filename=img_as_base64.txt")	
	io.WriteString(w, str_b64)
	
}




func base64_from_var(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	mode := r.FormValue("m")
	if mode == "" { mode = "1" }

	var str_src string
	if mode == "1"  { 
		str_src = conv.Img_jpeg_base64
	}
	if mode == "2"  { 
		str_src = conv.Img_rgba_base64
	}


	img,whichformat := conv.Base64_str_to_img(str_src)		
	c.Infof( "retrieved img from base64: format %v - type %T\n" , whichformat, img )

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "inline;filename=img_as_base64.txt")	
	io.WriteString(w, str_src)

	
}


func imageCache(w http.ResponseWriter, r *http.Request, dir , base string, c appengine.Context) {
	dsObj, _  := util.Buf_get(c , "util.WrapBlob_chart1")
	buff1, _  := conv.Vvbyte_to_string(dsObj.Vvbyte)
	
	img,whichformat := conv.Base64_str_to_img( buff1.String()  )	
	c.Infof( "retrieved img from base64: format %v - type %T\n" , whichformat, img )
	
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)		

}



func imageServe(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// prepare a cutout rect
	var p1, p2 image.Point
	p1.X, p1.Y = 10,10
	p2.X, p2.Y = 400,255
	var rec image.Rectangle = image.Rectangle{Min:p1,Max:p2}

	// prepare a line color
	lineCol := color.RGBA{}		
	lineCol.R, lineCol.G, lineCol.B = 255,244,22
	//lineCol.A = 0
	c.Infof( "brush color %#v \n" , lineCol)

	p := r.FormValue("p")
	if p == "" { p = "static/chartbg_400x960__480x1040__12x10.png" }	
	if p == "" { p = "static/pberg1.png" }	


	mode := r.FormValue("mode")
	if mode == "" { mode = "modified" }

	
	f, err := os.Open(p)
	if err != nil {  c.Errorf( "X3: %v \n" , err.Error() ); return  }
	defer f.Close()
	
	
	if mode == "direct" {
		c.Infof( "direct serving %v \n" , p )
		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, f)

	} else if mode == "unmodified" {
		img, whichformat, err := image.Decode(f)
		if err != nil {  c.Errorf( "X4: %v \n" , err.Error() ); return  }
		c.Infof( "serving unmodified - format %v - type %T\n" , whichformat, img )

		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, img, &jpeg.Options{Quality:jpeg.DefaultQuality})		

	} else if mode == "composited" {

		imgRGBA  := image.NewRGBA(rec)		
		off1 := imgRGBA.PixOffset(2,4)
		c.Infof( "offset is %v \n" , off1)
		for i := 0; i < 110; i++ {
			lineCol.A = uint8(i)
			imgRGBA.Set( i  , i,lineCol) 
			imgRGBA.Set( i+1, i,lineCol) 
			imgRGBA.Set( i+2, i,lineCol) 
			imgRGBA.Set( i+3, i,lineCol) 
		}		
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, imgRGBA, &jpeg.Options{Quality:jpeg.DefaultQuality})		

	} else if mode == "modified" {
		img, whichformat, err := image.Decode(f)
		if err != nil {  c.Errorf( "X5: %v \n" , err.Error() ); return  }
		c.Infof( "serving modified - format %v %T\n" , whichformat , img)
		
		
		switch t := img.(type) {

			default:
				c.Errorf("unexpected type %T", t)	   

			case *image.YCbCr:
				imgXFull,ok := img.(*image.YCbCr)
				util_err.Err_http(w,r,ok)
				
				imgXCutout := imgXFull.SubImage(rec)
				w.Header().Set("Content-Type", "image/jpeg")
				jpeg.Encode(w, imgXCutout, &jpeg.Options{Quality:jpeg.DefaultQuality})		


			case *image.RGBA:
				imgXFull,ok := img.(*image.RGBA)
				util_err.Err_http(w,r,ok)

				drawLineClosure := funcDrawLiner(lineCol,imgXFull)

				xb,yb := 40,440
				P0 := image.Point{xb +  0 ,yb -  0}
				drawLineClosure( P0, lineCol,imgXFull ) 

				for i := 0; i < 1; i++ {

					P1 := image.Point{xb + 80 ,yb - 80}
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +160; P1.Y = yb -160
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +240; P1.Y = yb -240
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +320; P1.Y = yb -320
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +400; P1.Y = yb -400
					drawLineClosure( P1, lineCol,imgXFull) 


					drawLineClosure = funcDrawLiner(lineCol,imgXFull)
					yb = 440
					P0 = image.Point{xb +  0 ,yb -  0}
					drawLineClosure( P0, lineCol,imgXFull ) 


					P1 = image.Point{xb + 80 ,yb - 40}
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +160; P1.Y = yb - 90
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +240; P1.Y = yb -120
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +320; P1.Y = yb -300
					drawLineClosure( P1, lineCol,imgXFull) 
					P1.X = xb +400; P1.Y = yb -310
					drawLineClosure( P1, lineCol,imgXFull) 

				}

				var imgXCutout *image.RGBA
				if true {
					imgXCutout = imgXFull
				} else {
					imgXCutout,ok = imgXFull.SubImage(rec).(*image.RGBA)
					util_err.Err_http(w,r,ok)
				}

				str_b64_img := conv.Rgba_img_to_base64_str( imgXCutout )
				//c.Infof("%v", str_b64_img[:100])

				vvbyte,_     := conv.String_to_vvbyte(str_b64_img)
				key_combi,_  := util.Buf_put(c , util.WrapBlob{"chart1",vvbyte} , "chart1" )

				dsObj,_  := util.Buf_get(c , key_combi)
				buff1,_  := conv.Vvbyte_to_string(dsObj.Vvbyte)

				img,whichformat := conv.Base64_str_to_img( buff1.String()  )	
				c.Infof( "retrieved img from base64: format %v - type %T\n" , whichformat, img )

				w.Header().Set("Content-Type", "image/png")
				png.Encode(w, img)
				
			// end case
				
		}

		

	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write( []byte("tell p=src&mode=<direct|unmodified|modified|composited>"))		
	}
	
	

	
}



func init() {
	http.HandleFunc("/image-analyze", imageAnalyze)	
	http.HandleFunc("/image-serve", imageServe)	

	http.HandleFunc("/base64-from-file", base64_from_file)	
	http.HandleFunc("/base64-from-var" , base64_from_var)	


	http.HandleFunc("/image-cache", adapterAddC(imageCache) )	
}
