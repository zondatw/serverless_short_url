import Cookies from 'cookies'

export default async (req, res) => {
  if (req.method === 'POST') {
    const cookies = new Cookies(req, res)
    const requestData = {
      email: req["body"]["email"],
      password: req["body"]["password"],
    }
    const response = await fetch(process.env.NEXTAUTH_URL, {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/json'
      }),
      body: JSON.stringify(requestData),
    })
    const jsonData = await response.json()
    // Set a cookie
    cookies.set('apikey', jsonData['idToken'], {
        httpOnly: true // true by default
    })
    res.status(200).json(jsonData)
  }
}
