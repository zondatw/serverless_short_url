import getConfig from 'next/config'

import Layout from '../components/layout'
import ShortUrlList from '../components/shortUrlList'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


export async function getServerSideProps() {
  var apiURL = process.env.apiURL || publicRuntimeConfig.apiURL
  return {
    props: {
      apiURL: apiURL
    }
  }
}

export default function Dashboard({apiURL}) {
  return (
    <Layout title="Dashboard">
      <div className="px-4 py-6 sm:px-0">
        <ShortUrlList apiUrl={apiURL} />
      </div>
    </Layout>
  )
}