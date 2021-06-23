import Router from 'next/router'

import Cookies from 'cookies'

import Layout from '../components/layout'
import LoginCard  from '../components/login'

// export async function getServerSideProps() {

//   return {
//     props: {
//     }
//   }
// }


export default function Login() {
  // console.log(`Login ${apikey}`)
  // if (apikey !== "") {
  //   Router.push("/")
  // }

  return (
    <Layout title="Login">
      <LoginCard />
    </Layout>
  )
}

// Login.getInitialProps = async ({ req, res }) => {
//   console.log(`init ${req}`)
//     const cookies = new Cookies(req, res)
//     const apikey = cookies.get('apikey')
//     console.log(`init ${apikey}`)

//     return {
//       apikey: apikey
//     }
// }