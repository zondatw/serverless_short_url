import { useState } from 'react';

export default function RegisterCard({ apiURL }) {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');

  function copyToClipboard(e) {
    navigator.clipboard.writeText(shortUrl)
  };

  function registerUrl() {
    setShortUrl("")
    var requestData = {
      url: url
    }
    fetch(apiUrl + "Register", {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/json'
      }),
      body: JSON.stringify(requestData),
    }).then((response) => {
      console.log(response)
      return response.json()
    }).then((jsonData) => {
      setShortUrl(jsonData["url"])
    }).catch(error => console.error('Error:', error))
  }

  return (
    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div className="flex flex-wrapbg-white py-8 px-10 text-center rounded-md shadow-lg max-w-lg mx-auto">
          <div className="flex flex-col">
            <div className="flex flex-row">
              <div>
                <input
                  type="text"
                  value={url}
                  onChange={e => {setUrl(e.currentTarget.value)}}
                  className="sm:text-xl border-2 border-blue-400 focus:ring-indigo-500 focus:border-indigo-500 flex-1 block w-full rounded-md  pt-3 pb-4 px-8 inline"
                  placeholder="URL"
                />
              </div>
              <div className="px-2">
                <button
                  onClick={() => registerUrl()}
                  className="rounded-md bg-gradient-to-r from-blue-400 to-indigo-500 text-xl text-white pt-3 pb-4 px-8 inline"
                >Register</button>
              </div>
            </div>

            {shortUrl? 
              <div>
                <div className="hidden sm:block" aria-hidden="true">
                  <div className="py-5">
                    <div className="border-t border-gray-200" />
                  </div>
                </div>
                <h1
                  onClick={copyToClipboard}
                >{shortUrl}</h1>
              </div>
              :
              <></>
            }
          </div>
        </div>
      </div>
    </div>
  )
}