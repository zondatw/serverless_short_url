import Cookies from 'cookies'
import getConfig from 'next/config'

import Layout from '../components/layout'
import { useAuthContext } from "../context/auth"
import DashboardComp from '../components/dashbaordComp'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


export async function getServerSideProps(context) {
  var apiURL = process.env.apiURL || publicRuntimeConfig.apiURL
  const cookies = new Cookies(context.req, context.res)
  var apikey = cookies.get('apikey')
  if (apikey == null) {
    apikey = ""
  }

  return {
    props: {
      apikey: apikey,
      isSignIn: apikey != "",
      apiURL: apiURL,
    }
  }
}

export default function Dashboard({apikey, isSignIn, apiURL}) {
  const {authSignIn, authSignOut} = useAuthContext()
  if (isSignIn) {
    authSignIn()
  } else {
    authSignOut()
  }

  return (
    <Layout title="Dashboard">
      <div className="px-4 py-6 sm:px-0">
        <DashboardComp apiUrl={apiURL} apikey={apikey} length={5} />
      </div>
    </Layout>
  )
}