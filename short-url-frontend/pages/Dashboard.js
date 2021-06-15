import Layout from '../components/layout'
import ShortUrlList from '../components/shortUrlList'

export default function Dashboard() {
  return (
    <Layout title="Dashboard">
      <div className="px-4 py-6 sm:px-0">
        <ShortUrlList />
      </div>
    </Layout>
  )
}