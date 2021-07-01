import Error from 'next/error'
import Cookies from 'cookies'
import getConfig from 'next/config'

import Layout from '../../components/layout'
import { useAuthContext } from "../../context/auth"
import ShortUrlDetail from '../../components/shorUrlDetail'

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
      errorCode: errorCode,
      hash: hash,
      data: data,
      isSignIn: apikey != "",
    }
  }
}


export default function Detail({errorCode, hash, data, isSignIn}) {
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
        data={data}
      />
    </Layout>
  )
}