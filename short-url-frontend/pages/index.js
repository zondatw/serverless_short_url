import Cookies from 'cookies'

import Layout from '../components/layout'
import { useAuthContext } from "../context/auth"


export async function getServerSideProps(context) {
  const cookies = new Cookies(context.req, context.res)
  const apikey = cookies.get('apikey')

  return {
    props: {
      isSignIn: apikey != null
    }
  }
}


export default function Home({isSignIn}) {
  const {authSignIn, authSignOut} = useAuthContext()
  if (isSignIn) {
    authSignIn()
  } else {
    authSignOut()
  }

  return (
    <Layout title="Home">
      <h1>Welcome!</h1>
    </Layout>
  )
}