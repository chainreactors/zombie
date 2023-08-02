package core

import (
	"fmt"
	"unsafe"
)

func CVAL(p *[]uint8) int {
	a := int((*p)[0])
	*p = (*p)[1:]
	return a
}

func CVAL2(p *[]uint8, v *uint16) {
	*v = *((*uint16)(unsafe.Pointer(&(*p)[0])))
	*p = (*p)[2:]
}

func CVAL3(p *[]uint8, v *[3]uint8) {
	(*v)[0] = (*p)[0]
	(*v)[1] = (*p)[1]
	(*v)[2] = (*p)[2]
	*p = (*p)[3:]
}

func REPEAT(f func(), count *int, x *int, width int) {
	for (*count & ^0x7) != 0 && ((*x + 8) < width) {
		for i := 0; i < 8; i++ {
			f()
			*count = *count - 1
			*x = *x + 1
		}
	}

	for (*count > 0) && (*x < width) {
		f()
		*count = *count - 1
		*x = *x + 1
	}
}

/* 1 byte bitmap decompress */
func decompress1(output *[]uint8, width, height int, input []uint8, size int) bool {
	var (
		prevline, line, count            int
		offset, code                     int
		x                                int = width
		opcode                           int
		lastopcode                       int8 = -1
		insertmix, bicolour, isfillormix bool
		mixmask, mask                    uint8
		colour1, colour2                 uint8
		mix                              uint8 = 0xff
		fom_mask                         uint8
	)
	out := *output
	for len(input) != 0 {
		fom_mask = 0
		code = CVAL(&input)
		opcode = code >> 4
		/* Handle different opcode forms */
		switch opcode {
		case 0xc, 0xd, 0xe:
			opcode -= 6
			count = int(code & 0xf)
			offset = 16
			break
		case 0xf:
			opcode = code & 0xf
			if opcode < 9 {
				count = int(CVAL(&input))
				count |= int(CVAL(&input) << 8)
			} else {
				count = 1
				if opcode < 0xb {
					count = 8
				}
			}
			offset = 0
			break
		default:
			opcode >>= 1
			count = int(code & 0x1f)
			offset = 32
			break
		}
		/* Handle strange cases for counts */
		if offset != 0 {
			isfillormix = ((opcode == 2) || (opcode == 7))
			if count == 0 {
				if isfillormix {
					count = int(CVAL(&input)) + 1
				} else {
					count = int(CVAL(&input) + offset)
				}
			} else if isfillormix {
				count <<= 3
			}
		}
		/* Read preliminary data */
		switch opcode {
		case 0: /* Fill */
			if (lastopcode == int8(opcode)) && !((x == width) && (prevline == 0)) {
				insertmix = true
			}
			break
		case 8: /* Bicolour */
			colour1 = uint8(CVAL(&input))
			colour2 = uint8(CVAL(&input))
			break
		case 3: /* Colour */
			colour2 = uint8(CVAL(&input))
			break
		case 6: /* SetMix/Mix */
			fallthrough
		case 7: /* SetMix/FillOrMix */
			mix = uint8(CVAL(&input))
			opcode -= 5
			break
		case 9: /* FillOrMix_1 */
			mask = 0x03
			opcode = 0x02
			fom_mask = 3
			break
		case 0x0a: /* FillOrMix_2 */
			mask = 0x05
			opcode = 0x02
			fom_mask = 5
			break
		}
		lastopcode = int8(opcode)
		mixmask = 0
		/* Output body */
		for count > 0 {
			if x >= width {
				if height <= 0 {
					return false
				}

				x = 0
				height--
				prevline = line
				line = height * width
			}
			switch opcode {
			case 0: /* Fill */
				if insertmix {
					if prevline == 0 {
						out[x+line] = mix
					} else {
						out[x+line] = out[prevline+x] ^ mix
					}
					insertmix = false
					count--
					x++
				}
				if prevline == 0 {
					REPEAT(func() {
						out[x+line] = 0
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[x+line] = out[prevline+x]
					}, &count, &x, width)
				}
				break
			case 1: /* Mix */
				if prevline == 0 {
					REPEAT(func() {
						out[x+line] = mix
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[x+line] = out[prevline+x] ^ mix
					}, &count, &x, width)
				}
				break
			case 2: /* Fill or Mix */
				if prevline == 0 {
					REPEAT(func() {
						mixmask <<= 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[x+line] = mix
						} else {
							out[x+line] = 0
						}
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						mixmask = mixmask << 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[x+line] = out[prevline+x] ^ mix
						} else {
							out[x+line] = out[prevline+x]
						}
					}, &count, &x, width)
				}
				break
			case 3: /* Colour */
				REPEAT(func() {
					out[x+line] = colour2
				}, &count, &x, width)
				break
			case 4: /* Copy */
				REPEAT(func() {
					out[x+line] = uint8(CVAL(&input))
				}, &count, &x, width)
				break
			case 8: /* Bicolour */
				REPEAT(func() {
					if bicolour {
						out[x+line] = colour2
						bicolour = false
					} else {
						out[x+line] = colour1
						bicolour = true
						count++
					}
				}, &count, &x, width)

				break

			case 0xd: /* White */
				REPEAT(func() {
					out[x+line] = 0xff
				}, &count, &x, width)
				break
			case 0xe: /* Black */
				REPEAT(func() {
					out[x+line] = 0
				}, &count, &x, width)
				break
			default:
				fmt.Printf("bitmap opcode 0x%x\n", opcode)
				return false
			}
		}
	}
	return true
}

/* 2 byte bitmap decompress */
func decompress2(output *[]uint8, width, height int, input []uint8, size int) bool {
	var (
		prevline, line, count            int
		offset, code                     int
		x                                int = width
		opcode                           int
		lastopcode                       int = -1
		insertmix, bicolour, isfillormix bool
		mixmask, mask                    uint8
		colour1, colour2                 uint16
		mix                              uint16 = 0xffff
		fom_mask                         uint8
	)

	out := make([]uint16, width*height)
	for len(input) != 0 {
		fom_mask = 0
		code = CVAL(&input)
		opcode = code >> 4
		/* Handle different opcode forms */
		switch opcode {
		case 0xc, 0xd, 0xe:
			opcode -= 6
			count = code & 0xf
			offset = 16
			break
		case 0xf:
			opcode = code & 0xf
			if opcode < 9 {
				count = CVAL(&input)
				count |= CVAL(&input) << 8
			} else {
				count = 1
				if opcode < 0xb {
					count = 8
				}
			}
			offset = 0
			break
		default:
			opcode >>= 1
			count = code & 0x1f
			offset = 32
			break
		}

		/* Handle strange cases for counts */
		if offset != 0 {
			isfillormix = ((opcode == 2) || (opcode == 7))
			if count == 0 {
				if isfillormix {
					count = CVAL(&input) + 1
				} else {
					count = CVAL(&input) + offset
				}
			} else if isfillormix {
				count <<= 3
			}
		}
		/* Read preliminary data */
		switch opcode {
		case 0: /* Fill */
			if (lastopcode == opcode) && !((x == width) && (prevline == 0)) {
				insertmix = true
			}
			break
		case 8: /* Bicolour */
			CVAL2(&input, &colour1)
			CVAL2(&input, &colour2)
			break
		case 3: /* Colour */
			CVAL2(&input, &colour2)
			break
		case 6: /* SetMix/Mix */
			fallthrough
		case 7: /* SetMix/FillOrMix */
			CVAL2(&input, &mix)
			opcode -= 5
			break
		case 9: /* FillOrMix_1 */
			mask = 0x03
			opcode = 0x02
			fom_mask = 3
			break
		case 0x0a: /* FillOrMix_2 */
			mask = 0x05
			opcode = 0x02
			fom_mask = 5
			break
		}
		lastopcode = opcode
		mixmask = 0
		/* Output body */
		for count > 0 {
			if x >= width {
				if height <= 0 {
					return false
				}

				x = 0
				height--
				prevline = line
				line = height * width
			}
			switch opcode {
			case 0: /* Fill */
				if insertmix {
					if prevline == 0 {
						out[x+line] = mix
					} else {
						out[x+line] = out[prevline+x] ^ mix
					}
					insertmix = false
					count--
					x++
				}
				if prevline == 0 {
					REPEAT(func() {
						out[x+line] = 0
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[x+line] = out[prevline+x]
					}, &count, &x, width)
				}
				break
			case 1: /* Mix */
				if prevline == 0 {
					REPEAT(func() {
						out[x+line] = mix
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[x+line] = out[prevline+x] ^ mix
					}, &count, &x, width)
				}
				break
			case 2: /* Fill or Mix */
				if prevline == 0 {
					REPEAT(func() {
						mixmask <<= 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[x+line] = mix
						} else {
							out[x+line] = 0
						}
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						mixmask = mixmask << 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[x+line] = out[prevline+x] ^ mix
						} else {
							out[x+line] = out[prevline+x]
						}
					}, &count, &x, width)
				}
				break
			case 3: /* Colour */
				REPEAT(func() {
					out[x+line] = colour2
				}, &count, &x, width)
				break
			case 4: /* Copy */
				REPEAT(func() {
					var a uint16
					CVAL2(&input, &a)
					out[x+line] = a
				}, &count, &x, width)

				break
			case 8: /* Bicolour */
				REPEAT(func() {
					if bicolour {
						out[x+line] = colour2
						bicolour = false
					} else {
						out[x+line] = colour1
						bicolour = true
						count++
					}
				}, &count, &x, width)

				break
			case 0xd: /* White */
				REPEAT(func() {
					out[x+line] = 0xffff
				}, &count, &x, width)
				break
			case 0xe: /* Black */
				REPEAT(func() {
					out[x+line] = 0
				}, &count, &x, width)
				break
			default:
				fmt.Printf("bitmap opcode 0x%x\n", opcode)
				return false
			}
		}
	}
	j := 0
	for _, v := range out {
		(*output)[j], (*output)[j+1] = PutUint16BE(v)
		j += 2
	}
	return true
}

// /* 3 byte bitmap decompress */
func decompress3(output *[]uint8, width, height int, input []uint8, size int) bool {
	var (
		prevline, line, count            int
		opcode, offset, code             int
		x                                int = width
		lastopcode                       int = -1
		insertmix, bicolour, isfillormix bool
		mixmask, mask                    uint8
		colour1                          = [3]uint8{0, 0, 0}
		colour2                          = [3]uint8{0, 0, 0}
		mix                              = [3]uint8{0xff, 0xff, 0xff}
		fom_mask                         uint8
	)
	out := *output
	for len(input) != 0 {
		fom_mask = 0
		code = CVAL(&input)
		opcode = code >> 4
		/* Handle different opcode forms */
		switch opcode {
		case 0xc, 0xd, 0xe:
			opcode -= 6
			count = code & 0xf
			offset = 16
			break
		case 0xf:
			opcode = code & 0xf
			if opcode < 9 {
				count = CVAL(&input)
				count |= CVAL(&input) << 8
			} else {
				count = 1
				if opcode < 0xb {
					count = 8
				}
			}
			offset = 0
			break
		default:
			opcode >>= 1
			count = code & 0x1f
			offset = 32
			break
		}

		/* Handle strange cases for counts */
		if offset != 0 {
			isfillormix = ((opcode == 2) || (opcode == 7))
			if count == 0 {
				if isfillormix {
					count = CVAL(&input) + 1
				} else {
					count = CVAL(&input) + offset
				}
			} else if isfillormix {
				count <<= 3
			}
		}
		/* Read preliminary data */
		switch opcode {
		case 0: /* Fill */
			if (lastopcode == opcode) && !((x == width) && (prevline == 0)) {
				insertmix = true
			}
			break
		case 8: /* Bicolour */
			CVAL3(&input, &colour1)
			CVAL3(&input, &colour2)
			break
		case 3: /* Colour */
			CVAL3(&input, &colour2)
			break
		case 6: /* SetMix/Mix */
			fallthrough
		case 7: /* SetMix/FillOrMix */
			CVAL3(&input, &mix)
			opcode -= 5
			break
		case 9: /* FillOrMix_1 */
			mask = 0x03
			opcode = 0x02
			fom_mask = 3
			break
		case 0x0a: /* FillOrMix_2 */
			mask = 0x05
			opcode = 0x02
			fom_mask = 5
			break
		}

		lastopcode = opcode
		mixmask = 0
		/* Output body */
		for count > 0 {
			if x >= width {
				if height <= 0 {
					return false
				}

				x = 0
				height--
				prevline = line
				line = height * width * 3
			}
			switch opcode {
			case 0: /* Fill */
				if insertmix {
					if prevline == 0 {
						out[3*x+line] = mix[0]
						out[3*x+line+1] = mix[1]
						out[3*x+line+2] = mix[2]
					} else {
						out[3*x+line] = out[prevline+3*x] ^ mix[0]
						out[3*x+line+1] = out[prevline+3*x+1] ^ mix[1]
						out[3*x+line+2] = out[prevline+3*x+2] ^ mix[2]
					}
					insertmix = false
					count--
					x++
				}
				if prevline == 0 {
					REPEAT(func() {
						out[3*x+line] = 0
						out[3*x+line+1] = 0
						out[3*x+line+2] = 0
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[3*x+line] = out[prevline+3*x]
						out[3*x+line+1] = out[prevline+3*x+1]
						out[3*x+line+2] = out[prevline+3*x+2]
					}, &count, &x, width)
				}
				break
			case 1: /* Mix */
				if prevline == 0 {
					REPEAT(func() {
						out[3*x+line] = mix[0]
						out[3*x+line+1] = mix[1]
						out[3*x+line+2] = mix[2]
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						out[3*x+line] = out[prevline+3*x] ^ mix[0]
						out[3*x+line+1] = out[prevline+3*x+1] ^ mix[1]
						out[3*x+line+2] = out[prevline+3*x+2] ^ mix[2]
					}, &count, &x, width)
				}
				break
			case 2: /* Fill or Mix */
				if prevline == 0 {
					REPEAT(func() {
						mixmask = mixmask << 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[3*x+line] = mix[0]
							out[3*x+line+1] = mix[1]
							out[3*x+line+2] = mix[2]
						} else {
							out[3*x+line] = 0
							out[3*x+line+1] = 0
							out[3*x+line+2] = 0
						}
					}, &count, &x, width)
				} else {
					REPEAT(func() {
						mixmask = mixmask << 1
						if mixmask == 0 {
							mask = fom_mask
							if fom_mask == 0 {
								mask = uint8(CVAL(&input))
								mixmask = 1
							}
						}
						if mask&mixmask != 0 {
							out[3*x+line] = out[prevline+3*x] ^ mix[0]
							out[3*x+line+1] = out[prevline+3*x+1] ^ mix[1]
							out[3*x+line+2] = out[prevline+3*x+2] ^ mix[2]
						} else {
							out[3*x+line] = out[prevline+3*x]
							out[3*x+line+1] = out[prevline+3*x+1]
							out[3*x+line+2] = out[prevline+3*x+2]
						}
					}, &count, &x, width)
				}
				break
			case 3: /* Colour */
				REPEAT(func() {
					out[3*x+line] = colour2[0]
					out[3*x+line+1] = colour2[1]
					out[3*x+line+2] = colour2[2]

				}, &count, &x, width)
				break
			case 4: /* Copy */
				REPEAT(func() {
					out[3*x+line] = uint8(CVAL(&input))
					out[3*x+line+1] = uint8(CVAL(&input))
					out[3*x+line+2] = uint8(CVAL(&input))
				}, &count, &x, width)
				break
			case 8: /* Bicolour */
				REPEAT(func() {
					if bicolour {
						out[3*x+line] = colour2[0]
						out[3*x+line+1] = colour2[1]
						out[3*x+line+2] = colour2[2]
						bicolour = false
					} else {
						out[3*x+line] = colour1[0]
						out[3*x+line+1] = colour1[1]
						out[3*x+line+2] = colour1[2]
						bicolour = true
						count++
					}
				}, &count, &x, width)
				break
			case 0xd: /* White */
				REPEAT(func() {
					out[3*x+line] = 0xff
					out[3*x+line+1] = 0xff
					out[3*x+line+2] = 0xff

				}, &count, &x, width)
				break
			case 0xe: /* Black */
				REPEAT(func() {
					out[3*x+line] = 0
					out[3*x+line+1] = 0
					out[3*x+line+2] = 0
				}, &count, &x, width)
				break
			default:
				fmt.Printf("bitmap opcode 0x%x\n", opcode)
				return false
			}
		}
	}

	return true
}

/* decompress a colour plane */
func processPlane(in *[]uint8, width, height int, output *[]uint8, j int) int {
	var (
		indexw   int
		indexh   int
		code     int
		collen   int
		replen   int
		color    uint8
		x        uint8
		revcode  int
		lastline int
		thisline int
	)
	ln := len(*in)

	lastline = 0
	indexh = 0
	i := 0
	for indexh < height {
		thisline = j + (width * height * 4) - ((indexh + 1) * width * 4)
		color = 0
		indexw = 0
		i = thisline

		if lastline == 0 {
			for indexw < width {
				code = CVAL(in)
				replen = int(code & 0xf)
				collen = int((code >> 4) & 0xf)
				revcode = (replen << 4) | collen
				if (revcode <= 47) && (revcode >= 16) {
					replen = revcode
					collen = 0
				}
				for collen > 0 {
					color = uint8(CVAL(in))
					(*output)[i] = uint8(color)
					i += 4

					indexw++
					collen--
				}
				for replen > 0 {
					(*output)[i] = uint8(color)
					i += 4
					indexw++
					replen--
				}
			}
		} else {
			for indexw < width {
				code = CVAL(in)
				replen = int(code & 0xf)
				collen = int((code >> 4) & 0xf)
				revcode = (replen << 4) | collen
				if (revcode <= 47) && (revcode >= 16) {
					replen = revcode
					collen = 0
				}
				for collen > 0 {
					x = uint8(CVAL(in))
					if x&1 != 0 {
						x = x >> 1
						x = x + 1
						color = -x
					} else {
						x = x >> 1
						color = x
					}
					x = (*output)[indexw*4+lastline] + color
					(*output)[i] = uint8(x)
					i += 4
					indexw++
					collen--
				}
				for replen > 0 {
					x = (*output)[indexw*4+lastline] + color
					(*output)[i] = uint8(x)
					i += 4
					indexw++
					replen--
				}
			}
		}
		indexh++
		lastline = thisline
	}
	return ln - len(*in)
}

/* 4 byte bitmap decompress */
func decompress4(output *[]uint8, width, height int, input []uint8, size int) bool {
	var (
		code             int
		onceBytes, total int
	)

	code = CVAL(&input)
	if code != 0x10 {
		return false
	}

	total = 1
	onceBytes = processPlane(&input, width, height, output, 3)
	total += onceBytes

	onceBytes = processPlane(&input, width, height, output, 2)
	total += onceBytes

	onceBytes = processPlane(&input, width, height, output, 1)
	total += onceBytes

	onceBytes = processPlane(&input, width, height, output, 0)
	total += onceBytes

	return size == total
}

/* main decompress function */
func Decompress(input []uint8, width, height int, Bpp int) []uint8 {
	size := width * height * Bpp
	output := make([]uint8, size)
	switch Bpp {
	case 1:
		decompress1(&output, width, height, input, size)
	case 2:
		decompress2(&output, width, height, input, size)
	case 3:
		decompress3(&output, width, height, input, size)
	case 4:
		decompress4(&output, width, height, input, size)
	default:
		fmt.Printf("Bpp %d\n", Bpp)
	}

	return output
}
