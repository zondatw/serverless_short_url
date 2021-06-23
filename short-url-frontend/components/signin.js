import { useState } from 'react'
import Router from 'next/router'

export default function SignInCard() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  function signin() {
    var signinApi = "/api/auth/signin"
    var requestData = {
      email: email,
      password: password,
      "returnSecureToken": true,
    }
    fetch(signinApi, {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/json'
      }),
      body: JSON.stringify(requestData),
    }).then((response) => {
      return response.json()
    }).then((jsonData) => {
      console.log(jsonData)
      if (!("error" in jsonData)) {
        Router.push('/')
      }
    }).catch(error => console.error('Error:', error))
  }

  return (
    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div className="flex flex-wrapbg-white py-8 px-10 text-left rounded-md shadow-lg max-w-sm mx-auto">
          <div className="flex flex-col">
            <div className="py-2">
              <h1 className="text-4xl font-black mb-4">Sign In</h1>
            </div>
            <div className="py-2">
              <input
                type="text"
                value={email}
                onChange={e => {setEmail(e.currentTarget.value)}}
                className="sm:text-lg border-2 border-blue-400 focus:ring-indigo-500 focus:border-indigo-500 flex-1 block w-full rounded-md pt-3 pb-4 px-8 inline"
                placeholder="email"
              />
            </div>
            <div className="py-2">
              <input
                type="password"
                value={password}
                onChange={e => {setPassword(e.currentTarget.value)}}
                className="sm:text-lg border-2 border-blue-400 focus:ring-indigo-500 focus:border-indigo-500 flex-1 block w-full rounded-md pt-3 pb-4 px-8 inline"
                placeholder="password"
              />
            </div>
            <div className="py-2">
              <button
                onClick={() => signin()}
                className="rounded-md bg-gradient-to-r from-blue-400 to-indigo-500 text-xl text-white pt-3 pb-4 px-8 inline"
              >Sign In</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}