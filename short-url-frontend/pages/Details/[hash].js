import getConfig from 'next/config'

import Layout from '../../components/layout'
import ShortUrlDetail from '../../components/shorUrlDetail'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


export async function getServerSideProps(context) {
  const { hash } = context.query
  var apiURL = process.env.apiURL || publicRuntimeConfig.apiURL

  const res = await fetch(`${apiURL}api/shorturl/${hash}`)
  const data = await res.json()

  return {
    props: {
      hash: hash,
      data: data,
    }
  }
}


export default function Detail({hash, data}) {

  return (
    <Layout title={`Detail ${hash}`}>
      <ShortUrlDetail
        hash={hash}
        data={data}
      />
    </Layout>
  )
}