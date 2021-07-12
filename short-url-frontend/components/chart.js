import React from 'react'
import {Line} from 'react-chartjs-2'

export default function LineChart({title, data}) {
  return (
    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div className="flex flex-wrapbg-white py-8 px-10 rounded-md shadow-md max-w-lg mx-auto">
          <div className="flex flex-col">
            <h2>{title}</h2>
            <Line
              data={data}
              width={400}
              height={400}
            />
          </div>
        </div>
      </div>
    </div>
  )
}