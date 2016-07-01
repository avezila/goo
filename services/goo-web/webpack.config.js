var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var devFlagPlugin = new webpack.DefinePlugin({
  __DEV__: JSON.stringify(JSON.parse(process.env.DEBUG || 'false'))
});
var precss       = require('precss');
var autoprefixer = require('autoprefixer');
 


module.exports = {
  devtool: 'eval',
  entry: [
    'webpack-dev-server/client?http://127.0.0.1:'+process.env.PORT,
    'webpack/hot/only-dev-server',
    './src/index'
  ],
  devServer: {
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PATCH, PUT, DELETE, OPTIONS", 
    }
  },
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'bundle.js',
    publicPath: '/static/'
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoErrorsPlugin(),
    devFlagPlugin,
    new ExtractTextPlugin('app.css'),
    new webpack.DefinePlugin({
      'process.env.SERVER': `"${process.env.SERVER}"`
    })
  ],
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        loaders: ['react-hot', 'babel'],
        include: path.join(__dirname, 'src')
      },
      { test: /\.css$/, loader: ExtractTextPlugin.extract('css-loader!postcss-loader') },
      { test: /\.(sass|scss)$/, loader: ExtractTextPlugin.extract('css-loader!postcss-loader!sass-loader') }
    ]
  },
  postcss: function () {
    return [precss, autoprefixer];
  },
  resolve: {
    extensions: ['', '.js', '.json']
  }
};
