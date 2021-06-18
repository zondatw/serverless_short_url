import useSWR from 'swr'
import { useState } from 'react';

import Pagination from './pagination'
import ShortUrlTable from './shortUrlTable';

var tempData = {
  "data": null,
  "start": "",
  "next": "",
}
var prevTable = {
  "": ""
}
const fetcher = (url, start, length) => fetch(url + `?start=${start}&length=${length}`).then(r => r.json());

/* Reference: https://tailwindui.com/components/application-ui/lists/tables */
export default function DashboardComp({ apiUrl, length }) {
  const [start, setStart] = useState("")
  const { data, error } = useSWR([apiUrl + 'api/shorturl/', start, length], fetcher)

  function toPrev() {
    setStart(prevTable[tempData["start"]])
  }
  function toNext() {
    setStart(tempData["next"])
  }

  if (data && !(data["start"] in prevTable)) prevTable[data["start"]] = tempData["start"]
  
  if (data) tempData = Object.assign(data)
  return (
    <div className="flex flex-col">
      <ShortUrlTable
        data={tempData["data"]}
      />
      <Pagination
        prev={toPrev}
        next={toNext}
      />
    </div>
  )
}