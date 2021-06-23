import Router from 'next/router'

import Cookies from 'cookies'

import Layout from '../components/layout'
import LoginCard  from '../components/login'


export default function Login() {
  return (
    <Layout title="Login">
      <LoginCard />
    </Layout>
  )
}