import useSWR from 'swr'

import LineChart from "./chart"

const fetcher = (url, year, month, ) => fetch(url + `?year=${year}&month=${month}`).then(r => r.json())

export default function DailyReportLineChart({title, apiUrl, hash}) {
  var { data, error } = useSWR([apiUrl + `api/shorturlreport/daily/${hash}`, 2021, 7], fetcher)

  var labels = []
  var counts = []
  if (data) {
    data["dates"].forEach(dailyData => {
      labels.push(dailyData["date"])
      counts.push(dailyData["count"])
      console.log(labels)
      console.log(counts)
    })
  }

  var charData  = {
    labels: labels,
    datasets: [
      {
        label: 'Count',
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
    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div className="flex flex-wrapbg-white py-8 px-10 rounded-md shadow-md max-w-4xl mx-auto">
          <div className="flex flex-col">
            <h2>{title}</h2>
            <LineChart data={charData} />
          </div>
        </div>
      </div>
    </div>
  )
}