import Layout from '../components/layout'
import ShortUrlList from '../components/shortUrlList'


export default function Dashboard() {
  const shortUrls = [
    {
      hash: "XXXXX",
      target: "https://google.com",
      type: "url",
    },
    {
      hash: "OOOOO",
      target: "https://yahoo.com.tw",
      type: "url",
    },
  ]

  return (
    <Layout title="Dashboard">
      <div className="px-4 py-6 sm:px-0">
        <ShortUrlList shortUrls={shortUrls} />
      </div>
    </Layout>
  )
}