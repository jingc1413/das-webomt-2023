
const { defineConfig } = require('@vue/cli-service');
const AutoImport = require('unplugin-auto-import/webpack');
const Components = require('unplugin-vue-components/webpack');

const IconsResolver = require('unplugin-icons/resolver');
const Icons = require('unplugin-icons/webpack');
const { ElementPlusResolver } = require('unplugin-vue-components/resolvers');
function resolve(dir) {
  return path.join(__dirname, dir)
}

const path = require('path')
module.exports = defineConfig({
  transpileDependencies: true,
  productionSourceMap: process.env.NODE_ENV === 'development',
  lintOnSave: true,
  devServer: {
    client: {
      overlay: false,
    },
  },
  configureWebpack: {
    resolve: {
      extensions: [".ts", ".tsx", ".js", ".json"],
      alias: {
        '@': resolve('src')
      }
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          loader: 'ts-loader',
          exclude: /node_modules/,
          options: {
            appendTsSuffixTo: [/\.vue$/],
          }
        }
      ]
    },
    plugins: [
      //配置webpack自动按需引入element-plus，
      require('unplugin-element-plus/webpack')({
        libs: [{
          libraryName: 'element-plus',
          esModule: true,
          resolveStyle: (name) => {
            return `element-plus/theme-chalk/${name}.css`
          },
        }, ]
      }),
      AutoImport({
        resolvers: [
          // 自动导入图标组件
          IconsResolver({
            prefix: 'Icon',
          }),
          ElementPlusResolver()
        ]
      }),
      Components({
        resolvers: [
          // 自动注册图标组件
          IconsResolver({
            enabledCollections: ['ep'],
          }),
          ElementPlusResolver()
        ]
      }),
      Icons({
        autoInstall: true,
      }),
    ],
  },
  chainWebpack(config) {
    // set svg-sprite-loader
    config.module
        .rule('svg')
        .exclude.add(resolve('src/icons'))
        .end()
    config.module
        .rule('icons')
        .test(/\.svg$/)
        .include.add(resolve('src/icons'))
        .end()
        .use('svg-sprite-loader')
        .loader('svg-sprite-loader')
        .options({
          symbolId: 'icon-[name]'
        })
        .end()
  },
})
