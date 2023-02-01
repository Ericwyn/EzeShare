#/usr/bin/env bash
function _myscript(){
    if [[ "${COMP_CWORD}" == "1" ]];
    then
        COMP_WORD="-receiver -sender"
        COMPREPLY=($(compgen -W "$COMP_WORD" -- ${COMP_WORDS[${COMP_CWORD}]}))
    else
        if [[ "${COMP_CWORD}" == "2" ]];
        then
            case ${COMP_WORDS[$[$COMP_CWORD-1]]} in
            -receiver)
            COMP_WORD_2="-ip -debug"
            echo ''
            echo '2222-receiver'
            echo ''
            COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
            ;;

            -sender)
            COMP_WORD_2="-ip -f -debug"
            COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
            ;;

            esac

        elif [[ "${COMP_CWORD}" == "3" ]];
        then
            case ${COMP_WORDS[$[$COMP_CWORD-1]]} in
            -f)
            COMPREPLY=($(compgen -f ${COMP_WORDS[${COMP_CWORD}]}))
            ;;

            esac
        fi
    fi
}
# 注册命令补全函数
complete -F _myscript ezeshare

# 使用
# sudo ln -s ./eze-share-completion.bash /usr/share/bash-completion/completions/EzeShare
# 添加到系统 completion 设置里
