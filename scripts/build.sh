#!/bin/bash

pwd=`pwd`
src="$pwd/../src"
output="$pwd/../output"

function print_usage() {
  echo "USAGE:"
  echo $'\t'"$0 Service ServiceName"

  echo "e.g.:"
  echo $'\t'"$0 all"
  echo $'\t'"$0 task srv"
  echo $'\t'"$0 notify worker"
  echo $'\t'"$0 auth all"

  exit 1;
}

function build() {
  # 构建
  # $1 服务路径
  # $2 服务类型路径

  echo "[I] Prepare building: $1-$2, output: $output"

  # 无法构建没有main.go文件的服务路径
  if [ ! -f "$src/$1/$2/main.go" ]; then
    echo "[E] Failed to build $1-$2, no main.go under dir."
    exit 2
  fi

  # 判断构建结果路径是否存在，不存在就创建
  if [ ! -d "$output" ]; then
    mkdir -p "$output"
  fi

  # 切换目录经并开始构建
  cd "$src/$1/$2"
  CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o "$output/$1-$2"

  # 构建是否成功
  if [ $? -eq 0 ]; then
    echo "[I] Success build $1-$2"
  else
    echo "[E] Failed build $1-$2"
  fi
}

function traversal_build() {
  # 遍历目录
  # $1 基础目录，可选

  if [ $# -gt 0 ] ; then
    base_dir="$src/$1"
  else
    base_dir=$src
  fi

  # 列出基础目录下所有文件和目录进行遍历
  for layer1 in `ls $base_dir`
  do
    # 如果当文件是目录就继续
    if [ -d "$base_dir/$layer1" ]; then
      # 如果传入了基础目录，就只遍历1层
      if [ $# -gt 0 ]; then
        # 遇到目录下有main.go文件的的就认为其是服务的目录并构建
        if [ -f "$base_dir/$layer1/main.go" ]; then
          build $1 $layer1
        fi
        continue
      fi

      # 没有传入基础目录，遍历2层
      for layer2 in `ls "$base_dir/$layer1"`
      do
        # 遇到目录下有main.go文件的的就认为其是服务的目录并构建
        if [ -d "$base_dir/$layer1/$layer2" ] && [ -f "$base_dir/$layer1/$layer2/main.go" ]; then
          build $layer1 $layer2
        fi
      done
    fi
  done
}

if [ $# -lt 1 ] ; then
  print_usage
fi

if [ $# -eq 2 ]; then
  if [ $2 == "all" ]; then
    traversal_build $1
  else
    build $1 $2
  fi

  echo "[I] Build finished"
  exit 0
fi

if [ $# -eq 1 ]; then
  if [ $1 == "all" ]; then
    traversal_build
  else
    print_usage
  fi

  echo "[I] Build finished"
  exit 0
fi