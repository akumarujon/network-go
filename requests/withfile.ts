const data = new FormData()


const file = await Deno.readFile(Deno.cwd() + "/assets"  + '/default.jpg')
data.append('picture', new Blob([file]), 'pfp.jpg')
data.append("username", "test")
data.append("password", "test")
data.append("email", "akumaru.senju@gmail.com")

const response = await fetch('https://81d7-185-213-230-13.ngrok-free.app/signup', {
    method: 'POST',
    body: data
})

console.log(response.status)
console.log(await response.json())

