package gl

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)


type ShaderInfo struct {
	ShaderType uint32
	Filename   string
	Shader     uint32
}

func readShader(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data) + "\x00", nil
}

func LoadShaders(shaders []ShaderInfo) (uint32, error) {
	if len(shaders) == 0 {
		return 0, nil
	}

	program := gl.CreateProgram()

	for _, entry := range shaders {
		shader := gl.CreateShader(entry.ShaderType)
		entry.Shader = shader

		source, err := readShader(entry.Filename)
		if err != nil {
			for _, e := range shaders {
				gl.DeleteProgram(e.Shader)
			}
			return 0, err
		}

		cstr, free := gl.Strs(source)
		gl.ShaderSource(shader, 1, cstr, nil)
		free()
		gl.CompileShader(shader)

		
		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

        if status == gl.FALSE {
            var len int32
            gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &len)

			log := strings.Repeat("\x00", int(len + 1))
            gl.GetShaderInfoLog(shader, len, nil, gl.Str(log))
            
            return 0, fmt.Errorf("shader compilation failed %v: %v", source, log)
        }

        gl.AttachShader(program, shader)
        
	}

	gl.LinkProgram(program)

    var linked int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &linked)
    if linked== gl.FALSE {
		var len int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &len)

		log := strings.Repeat("\x00", int(len + 1))
		gl.GetProgramInfoLog(program, len, nil, gl.Str(log))
		
		for _, e := range shaders {
			gl.DeleteShader(e.Shader)
			e.Shader = 0
		}

		return 0, fmt.Errorf("shader linking failed %v", log)
    }

    return program, nil
}