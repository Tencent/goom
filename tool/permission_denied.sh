#! /bin/bash

# 此脚本用于解决 permission denied 问题
# 执行脚本替换 sdk/pkg/tool/link 文件
# 在 link 文件中添加更改权限 (max_prot=7) 的脚本

############################################################
# Process the input options. Add options as needed.        #
############################################################

Help() {
  # Display Help
  echo "###############################################"
  echo "  __  __  __             __  __      ___ __  "
  echo " / _\`/  \/  \|\/|   |\/|/  \/  \`|__/|__ |__) "
  echo " \__>\__/\__/|  |   |  |\__/\__,|  \|___|  \ "
  echo "# 此脚本用于解决 MAC 环境下 permission denied 问题."
  echo "# 执行脚本替换 sdk/pkg/tool/link 文件, "
  echo "# 在 link 文件中添加更改编译产物执行权限(max_prot=7)的脚本."
  echo "使用方法: 选项 [-h|i|u|c]"
  echo "选项:"
  echo " -h     显示帮助."
  echo " -i     安装."
  echo " -u     卸载."
  echo " -c     查看安装状态."
  echo
  Check
}

Check() {
  echo `go version`
  echo "go root $(go env GOROOT)"

  TOOL_DIR=$(go env GOTOOLDIR)
  LINE=$(head -n 1 ${TOOL_DIR}/link)
  if [[ "$LINE" == "#!/usr/bin/env python3" ]]; then
    echo "当前状态: 已安装"
  else
    echo "当前状态: 未安装"
  fi
}

Install() {
  TOOL_DIR=$(go env GOTOOLDIR)
  WORK_DIR=$(cd $(dirname $0); pwd)
  LINE=$(head -n 1 ${TOOL_DIR}/link)

  if [[ "$LINE" == "#!/usr/bin/env python3" ]]; then
    echo "already installed."
    return 0
  fi

  if [ -e ${WORK_DIR}/link ]
  then
      mv ${TOOL_DIR}/link ${TOOL_DIR}/original_link
      cp ${WORK_DIR}/link ${TOOL_DIR}/link
      echo "replaced ${TOOL_DIR}/link with ${WORK_DIR}/link"
      echo "install success."
      return 0
  else
      echo "install fail: file not exists: ${TOOL_DIR}/link"
      return 1
  fi
}

Uninstall() {
  TOOL_DIR=$(go env GOTOOLDIR)
  LINE=$(head -n 1 ${TOOL_DIR}/link)

  if [[ "$LINE" == "#!/usr/bin/env python3" ]]; then
    mv ${TOOL_DIR}/original_link ${TOOL_DIR}/link
    echo "recovered ${TOOL_DIR}/link success."
  else
    echo "already recovered."
  fi
}

# Get the options
while getopts ":hiuc" option; do
  case $option in
  h) # display Help
    Help
    exit
    ;;
  i) # install
    Install
    exit
    ;;
  u) # uninstall
    Uninstall
    exit
    ;;
  c) # check
    Check
    exit
    ;;
  \?) # Invalid option
    echo "Error: Invalid option. use -h show help"
    Help
    exit
    ;;
  esac
done

Help
