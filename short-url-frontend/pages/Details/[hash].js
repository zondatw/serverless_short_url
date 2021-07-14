import Error from 'next/error'
import Cookies from 'cookies'
import getConfig from 'next/config'

import Layout from '../../components/layout'
import { useAuthContext } from "../../context/auth"
import ShortUrlDetail from '../../components/shorUrlDetail'
import DailyReportLineChart from '../../components/dailyReport'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()

export async function getServerSideProps(context) {
  var errorCode = false
  const { hash } = context.query
  var apiURL = process.env.apiURL || publicRuntimeConfig.apiURL
  const cookies = new Cookies(context.req, context.res)
  var apikey = cookies.get('apikey')
  var headers = {}
  if (apikey != null) {
    headers["Authorization"] = apikey
  } else {
    apikey = ""
  }

  const res = await fetch(`${apiURL}api/shorturl/${hash}`, {headers: headers})
  const data = await res.json()

  if ("error" in data) {
    errorCode = 403;
  }

  return {
    props: {
      apiURL: apiURL,
      apikey: apikey,
      errorCode: errorCode,
      hash: hash,
      shortUrlData: data,
      isSignIn: apikey != "",
    }
  }
}


export default function Detail({errorCode, hash, shortUrlData, isSignIn, apiURL, apikey}) {
  if (errorCode) {
    return <Error statusCode={errorCode} />
  }

  const {authSignIn, authSignOut} = useAuthContext()
  if (isSignIn) {
    authSignIn()
  } else {
    authSignOut()
  }

  return (
    <Layout title={`Detail ${hash}`}>
      <ShortUrlDetail
        hash={hash}
        data={shortUrlData}
      />
      <DailyReportLineChart title="Daily report" apiUrl={apiURL} apikey={apikey} hash={hash} />
    </Layout>
  )
}