import React from 'react'
import {Line} from 'react-chartjs-2'

export default function LineChart({data}) {
  return (
    <div>
      <h2>Line Example</h2>
      <Line
        data={data}
        width={400}
        height={400}
      />
    </div>
  )
}