<template>
  <div id="app">
    <h2>Modbus Manager</h2>
    <ul v-for="(reg) in registers_datas" :key="reg">
      <li>
        <h6>Register: {{reg[0]}}</h6>
        <h6>Value: {{reg[1]}}</h6>
        <button v-if="reg[1]==0" @click="action(reg[0], 1)">start</button>
        <button v-if="reg[1]==1" @click="action(reg[0], 0)">stop</button>
      </li>
    </ul>
  </div>
</template>

<script>

export default {
  name: 'App',
  data: function () {
    return {
      registers_datas: []
    }
  },
  created: function() {
    setInterval(() => {
      this.datas()
    }, 4000);
  },
  methods: {
    datas: function() {
      fetch("/datas")
      .then(resp => resp.json())
      .then(data => {
        for (const key in data) {
          var values = [key, data[key]]
        }
        this.registers_datas.push(values)
      })
      .catch(resp => console.error(resp))
    },
    action: function(register, set) {
      var res = {
        register: register,
        set: set
      }
      fetch("/action", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(res)
      })
      .then(resp => console.error(resp))
      .catch(resp => console.error(resp))
      this.datas()
    }
  }
}  
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
