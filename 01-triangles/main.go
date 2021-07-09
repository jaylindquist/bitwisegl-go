package main

import (
	"fmt"
	"log"
	"runtime"

	bgl "bitwiseor.com/bitwisegl/gl"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Triangles", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	initialize()

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

const numVertices = 6
var vbo uint32

func initialize() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertices := [6*2]float32{
         -0.90, -0.90,   0.85, -0.90,  -0.90,  0.85,
          0.90, -0.85,   0.90,  0.90,  -0.85,  0.90, 
    }
	
    gl.CreateBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferStorage(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(&vertices[0]), 0)

	shaders := []bgl.ShaderInfo{
		{ShaderType: gl.VERTEX_SHADER, Filename: "shaders/triangle.vert"},
		{ShaderType: gl.FRAGMENT_SHADER, Filename: "shaders/triangle.frag"},
	}

	program, err := bgl.LoadShaders(shaders)
	if err != nil{
		log.Fatalf("unable to load shaders %v: %v", shaders, err)
	}

    gl.UseProgram(program)

	vPosition := uint32(gl.GetAttribLocation(program, gl.Str("vPosition\x00")))
    gl.VertexAttribPointerWithOffset(vPosition, 2, gl.FLOAT, false, 2*4, 0)
    gl.EnableVertexAttribArray(vPosition)
}

func display() {
	black := [4]float32{0.0, 0.0, 0.0, 0.0}

	gl.ClearBufferfv(gl.COLOR, 0, &black[0])
	gl.BindVertexArray(vbo)
	gl.DrawArrays(gl.TRIANGLES, 0,numVertices)
}
