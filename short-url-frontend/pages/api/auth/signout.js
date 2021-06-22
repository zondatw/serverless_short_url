import Cookies from 'cookies'

export default async (req, res) => {
  const cookies = new Cookies(req, res)
  // Delete a cookie
  cookies.set('apikey')
  res.status(200).json({ status: "ok" })
}
