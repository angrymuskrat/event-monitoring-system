import React, { Suspense, lazy } from 'react'
import { Redirect, Route, Switch } from 'react-router-dom'

import { GlobalStyle } from './globalStyles'

// components
import Loading from '../../components/Loading/Loading'

// containers
import Navigation from '../Navigation'

// routes
const MapPage = lazy(() => import('../MapPage'))
const AboutPage = lazy(() => import('../../components/Pages/AboutPage.jsx'))
const TeamPage = lazy(() => import('../../components/Pages/TeamPage.jsx'))
const StartingPage = lazy(() => import('../StartingPage'))
const LoginPage = lazy(() => import('../LoginPage'))

function App({ session, auth }) {
  document.title = 'Event monitoring'
  return (
    <Suspense fallback={<Loading />}>
      <GlobalStyle />
      {auth ? (
        <>
          <Navigation />
          <Switch>
            <Route exact path="/" component={StartingPage} />
            <Route exact path="/login" component={LoginPage} />
            <Route
              exact
              path="/map/:city/:topLeft/:botRight/:time"
              component={MapPage}
            />
            <Route exact path="/about" component={AboutPage} />
            <Route exact path="/team" component={TeamPage} />
            <Route render={() => <Redirect to="/" />} />
          </Switch>
        </>
      ) : (
        <Switch>
          <Route exact path="/login" component={LoginPage} />
          <Route render={() => <Redirect to="/login" />} />
        </Switch>
      )}
    </Suspense>
  )
}

export default App
