import http from "node:http"


const server = http.createServer((req, res) => {
    res.end("Hello, World")
})


server.listen(3000, () => {
    console.log("listening on http://localhost:3000")
})