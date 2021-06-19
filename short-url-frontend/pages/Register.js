import getConfig from 'next/config'

import Layout from '../components/layout'
import RegisterCard from '../components/register'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


export async function getServerSideProps() {
  var apiURL = process.env.cloudFunctionUrl || publicRuntimeConfig.cloudFunctionUrl
  return {
    props: {
      apiURL: apiURL
    }
  }
}

export default function Register({apiURL}) {
  return (
    <Layout title="Register">
      <RegisterCard apiURL={apiURL} />
    </Layout>
  )
}