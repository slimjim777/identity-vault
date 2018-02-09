/*
 * Copyright (C) 2017-2018 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
import React, { Component } from 'react';
import Header from './components/Header'
import Footer from './components/Footer'
import Index from './components/Index'
import ModelList from './components/ModelList'
import ModelEdit from './components/ModelEdit'
import KeypairAdd from './components/KeypairAdd'
import KeypairGenerate from './components/KeypairGenerate'
import KeypairStore from './components/KeypairStore'
import AccountList from './components/AccountList'
import AccountForm from './components/AccountForm'
import AccountEdit from './components/AccountEdit'
import AccountDetail from './components/AccountDetail'
import AccountKeyForm from './components/AccountKeyForm'
import Keypair from './components/Keypair'
import SigningLog from './components/SigningLog'
import SystemUserForm from './components/SystemUserForm'
import UserList from './components/UserList'
import UserEdit from './components/UserEdit'
import Accounts from './models/accounts'
import Keypairs from './models/keypairs'
import Models from './models/models'
import {getAccount, saveAccount, sectionFromPath, sectionIdFromPath, subSectionIdFromPath} from './components/Utils'
import createHistory from 'history/createBrowserHistory'
import './sass/App.css'

const history = createHistory()

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      location: history.location,
      accounts: [],
      selectedAccount: getAccount() || {},
      keypairs: [],
      models: [],
    }

    history.listen(this.handleNavigation.bind(this))
    this.getAccounts()
  }

  handleNavigation(location) {
    this.setState({ location: location })
    window.scrollTo(0, 0)
  }

  getAccounts() {
    Accounts.list().then((response) => {
        var data = JSON.parse(response.body);
        var message = "";
        if (!data.success) {
            message = data.message;
        }

        var selectedAccount = this.state.selectedAccount;
        if ((!this.state.selectedAccount.code) && (!getAccount().code)) {
          // Set to the first in the account list
          if (data.accounts.length > 0) {
            selectedAccount = data.accounts[0]
            saveAccount(selectedAccount)
          }
        }

        this.setState({accounts: data.accounts, message: message});
        this.updateDataForRoute(selectedAccount)
    });
  }

  getKeypairs(account) {
    Keypairs.list(account).then((response) => {
        var data = JSON.parse(response.body);
        var message = "";
        if (!data.success) {
            message = data.message;
        }
        this.setState({keypairs: data.keypairs, message: message});
    });
  }

  getModels(account) {
    Models.list(account).then((response) => {
        var data = JSON.parse(response.body);
        var message = "";
        if (!data.success) {
            message = data.message;
        }
        this.setState({models: data.models, message: message});
    });
  }

  updateDataForRoute(selectedAccount) {
    var currentSection = sectionFromPath(window.location.pathname);

    if (currentSection === 'accounts') {
      this.getKeypairs(selectedAccount.AuthorityID)
      this.getModels(selectedAccount.AuthorityID)
    }
  }

  handleAccountChange = (account) => {
    console.log("===", account)
    saveAccount(account)
    this.setState({selectedAccount: account})

    this.updateDataForRoute(account)
  }

  renderModels() {
    const id = sectionIdFromPath(window.location.pathname, 'models')

    switch(id) {
      case 'new':
        return <ModelEdit token={this.props.token} />
      case '':
        return <ModelList token={this.props.token} />
      default:
        return <ModelEdit token={this.props.token} id={id} />
    }
  }

  renderAccounts() {
    const id = sectionIdFromPath(window.location.pathname, 'accounts')
    const sub = subSectionIdFromPath(window.location.pathname, 'accounts')

    // Only show the selected account
    var accounts = this.state.accounts.filter((a) => {
      return a.ID === this.state.selectedAccount.ID
    })

    switch(id) {
      case 'upload':
        return <AccountForm token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
      case 'new':
        return <AccountEdit token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
      case 'account':
        return <AccountEdit id={sub} token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
      case 'view':
        return <AccountDetail id={sub} token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
      case 'key-assertion':
        return <AccountKeyForm token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
      default:
        return <AccountList token={this.props.token} selectedAccount={this.state.selectedAccount} accounts={accounts} keypairs={this.state.keypairs} models={this.state.models} />
    }
  }

  renderUsers() {
    const id = sectionIdFromPath(window.location.pathname, 'users')

    switch(id) {
      case 'new':
        return <UserEdit token={this.props.token} />
      case '':
        return <UserList token={this.props.token} />
      default:
        return <UserEdit token={this.props.token} id={id} />
    }
  }

  renderKeypairs() {
    const id = sectionIdFromPath(window.location.pathname, 'signing-keys')

    switch(id) {
      case 'generate':
        return <KeypairGenerate token={this.props.token} />
      case 'new':
        return <KeypairAdd token={this.props.token} />
      case 'store':
        return <KeypairStore token={this.props.token} />
      default:
        return <Keypair token={this.props.token} />
    }
  }

  render() {

    var currentSection = sectionFromPath(window.location.pathname);
    console.log('---selected', this.state.selectedAccount)

    return (
      <div className="App">
          <Header token={this.props.token} accounts={this.state.accounts} selectedAccount={this.state.selectedAccount}
                  onAccountChange={this.handleAccountChange} />

          <div className="spacer" />
  
          {currentSection==='home'? <Index token={this.props.token} /> : ''}
          {currentSection==='notfound'? <Index token={this.props.token} error={true} /> : ''}

          {currentSection==='signing-keys'? this.renderKeypairs(): ''}
          {currentSection==='models'? this.renderModels() : ''}

          {currentSection==='accounts'? this.renderAccounts() : ''}
          {currentSection==='signinglog'? <SigningLog token={this.props.token} /> : ''}

          {currentSection==='systemuser'? <SystemUserForm token={this.props.token} /> : ''}

          {currentSection==='users'? this.renderUsers() : ''}

          <Footer />
      </div>
    )
  }
}

export default App;