<template>
  <div id="app">
    <div class="small-header-float-right">
      <el-select size="mini" v-model="lang" @change="languageChanged">
        <el-option label="English" value="en" />
        <el-option label="简体中文" value="zh" />
      </el-select>
    </div>
    <Home v-if="userInfo" :logoutMethod="logout"/>
    <Login v-else :loginMethod="login" :signupMethod="signup"/>
  </div>
</template>

<script>
import Home from './components/Home.vue'
import Login from './components/Login.vue'
import Config from './utils/config'
import i18n from './utils/i18n'
import md5 from 'md5'
import axios from 'axios'
import { getUserInfo, setUserInfo, removeUserInfo, getToken } from './utils/utils'

export default {
  name: 'App',
  components: {
    Home,
    Login
  },
  data() {
    return {
      userInfo: null,
      lang: 'en'
    }
  },
  beforeMount() {
    this.userInfo = getUserInfo()
    axios.defaults.headers.common['token'] = getToken()
  },
  methods: {
    logout() {
      removeUserInfo()
      axios.post('/api/logout', {
        username: this.userInfo.username,
      })
      .then()
      .catch(err => {
        console.log(err)
      })
      this.userInfo = null
    },
    login(username, password) {
      axios.post('/api/login', {
        username: username,
        password: md5(password)
      })
      .then(res => {
        this.userInfo = res.data
        setUserInfo(this.userInfo)
        axios.defaults.headers.common['token'] = this.userInfo.token
      })
      .catch(err => {
        this.$message({
          type: 'error',
          message: Config.lang.loginError
        })
        console.log(err)
      })
    },
    signup(username, password, cb) {
      axios.post('/api/signup', {
        username: username,
        password: md5(password)
      })
      .then(() => {
        cb()
      })
      .catch(err => {
        this.$message({
          type: 'error',
          message: Config.lang.signupError
        })
        console.log(err)
      })
    },
    languageChanged(val) {
      var lang = i18n[val]
      if (lang) {
        Config.lang = lang
      } else {
        console.log("invalid language selected:", val)
      }
    }
  }
}
</script>

<style>
#app {
  font-family: europa, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  margin: 8px 8px;
}

#app .small-header-float-right {
  width: 100%;
  margin-bottom: 8px;
  display: block;
}

</style>
