#include "stdio.h"


// big_endian: 1 
// little_endian: 2
int IsLittleEndian() {
    union {
        short value;
        char array[2];
    } u;
    u.value = 0x0102;
    if (u.array[0] == 1 && u.array[1] == 2){
        return 1;
    }else if (u.array[0] == 2 && u.array[1] == 1){
        return 2;
    }
    return -1;
}

int main() {
    
    int res;
    res = IsLittleEndian();
    printf("result is %d\n",res);
    if (res == 1) {
        printf("it is big endian");
    }
    if (res == 2){
        printf("it is little endian");
    }
    return 0;
}