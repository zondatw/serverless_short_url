import Error from 'next/error'
import Cookies from 'cookies'
import getConfig from 'next/config'
import useSWR from 'swr'

import Layout from '../../components/layout'
import { useAuthContext } from "../../context/auth"
import ShortUrlDetail from '../../components/shorUrlDetail'
import LineChart from '../../components/chart'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()


const fetcher = (url, year, month, ) => fetch(url + `?year=${year}&month=${month}`).then(r => r.json());

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
      apiURL,
      errorCode: errorCode,
      hash: hash,
      shortUrlData: data,
      isSignIn: apikey != "",
    }
  }
}


export default function Detail({errorCode, hash, shortUrlData, isSignIn, apiURL}) {
  if (errorCode) {
    return <Error statusCode={errorCode} />
  }

  const {authSignIn, authSignOut} = useAuthContext()
  if (isSignIn) {
    authSignIn()
  } else {
    authSignOut()
  }

  var { data, error } = useSWR([apiURL + `api/shorturlreport/daily/${hash}`, 2021, 7], fetcher)

  var labels = []
  var counts = []
  if (data) {
    data["dates"].forEach(dailyData => {
      labels.push(dailyData["date"])
      counts.push(dailyData["count"])
    })
  }

  var charData  = {
    labels: labels,
    datasets: [
      {
        label: 'My First dataset',
        fill: false,
        lineTension: 0.1,
        backgroundColor: 'rgba(75,192,192,0.4)',
        borderColor: 'rgba(75,192,192,1)',
        borderCapStyle: 'butt',
        borderDash: [],
        borderDashOffset: 0.0,
        borderJoinStyle: 'miter',
        pointBorderColor: 'rgba(75,192,192,1)',
        pointBackgroundColor: '#fff',
        pointBorderWidth: 1,
        pointHoverRadius: 5,
        pointHoverBackgroundColor: 'rgba(75,192,192,1)',
        pointHoverBorderColor: 'rgba(220,220,220,1)',
        pointHoverBorderWidth: 2,
        pointRadius: 1,
        pointHitRadius: 10,
        data: counts
      }
    ]
  };

  return (
    <Layout title={`Detail ${hash}`}>
      <ShortUrlDetail
        hash={hash}
        data={shortUrlData}
      />
      <LineChart data={charData} />
    </Layout>
  )
}