import React from 'react'
import {Line} from 'react-chartjs-2'

export default function LineChart({data}) {
  return (
    <Line
      data={data}
      width={800}
      height={400}
    />
  )
}