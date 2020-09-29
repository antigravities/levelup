const path = require("path");
const webpack = require("webpack");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CleanWebpackPlugin = require("clean-webpack-plugin").CleanWebpackPlugin;

module.exports = {
    entry: "./index.js",
    output: {
        path: path.resolve(__dirname + "/../static/"),
        filename: "kotoha.[contenthash].js"
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader']
            },
            {
                test: /\.html$/i,
                use: ['html-loader']
            },
            {
                test: /\.(png|svg|gif|jpg)$/i,
                use: ['file-loader']
            }
        ]
    },
    plugins: [
        new CleanWebpackPlugin(),
        new HtmlWebpackPlugin({
            template: "src/index.html",
        })
    ],
    mode: "production"
}