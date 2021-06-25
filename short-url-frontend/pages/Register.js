import Cookies from 'cookies'
import getConfig from 'next/config'

import Layout from '../components/layout'
import { useAuthContext } from "../context/auth"
import RegisterCard from '../components/register'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


export async function getServerSideProps(context) {
  var apiURL = process.env.cloudFunctionUrl || publicRuntimeConfig.cloudFunctionUrl
  const cookies = new Cookies(context.req, context.res)
  var apikey = cookies.get('apikey')
  if (apikey == null) {
    apikey = ""
  }

  return {
    props: {
      apikey: apikey,
      isSignIn: apikey != "",
      apiURL: apiURL
    }
  }
}

export default function Register({apikey, isSignIn, apiURL}) {
  const {authSignIn, authSignOut} = useAuthContext()
  if (isSignIn) {
    authSignIn()
  } else {
    authSignOut()
  }

  return (
    <Layout title="Register">
      <RegisterCard apiURL={apiURL} apikey={apikey} />
    </Layout>
  )
}