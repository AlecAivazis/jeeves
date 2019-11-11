task "generate" {
    description = "generate any files that need to created"
    command = "go generate ./..."
}

task "build" {
    command = "go build ."
}