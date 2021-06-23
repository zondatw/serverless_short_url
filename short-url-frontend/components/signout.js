import Router from 'next/router'

export default function SignOutCard() {

  function signout() {
    var signoutApi = "/api/auth/signout"
    fetch(signoutApi, {
      method: 'GET',
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
              <h1 className="text-4xl font-black mb-4">Are you sure?</h1>
            </div>
            <div className="py-2">
              <button
                onClick={() => signout()}
                className="rounded-md bg-gradient-to-r from-blue-400 to-indigo-500 text-xl text-white pt-3 pb-4 px-8 inline"
              >Sign Out</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}