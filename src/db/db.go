package db

import "C"

import (
  "github.com/jmhodges/levigo"
  "log"
  "strconv"
)

func CreateDB(pro string, ver string, address string, info string, break_pad_info string) {
  opts := levigo.NewOptions()
  defer opts.Close()
  opts.SetCache(levigo.NewLRUCache(3 << 30))
  opts.SetCreateIfMissing(true)

  ro := levigo.NewReadOptions()
  defer ro.Close()
  wo := levigo.NewWriteOptions()
  defer wo.Close()

  db_version, _ := levigo.Open("./"+pro+"/version_id.db", opts)
  defer db_version.Close()
  version, _ := db_version.Get(ro, []byte(ver))
  if version == nil {
    db_version.Put(wo, []byte(ver), []byte("0"))
  }

  db_id, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_id.db", opts)
  db_info, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_info.db", opts)
  db_breakpad_info, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", opts)
  db_count, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_count.db", opts)
  defer db_id.Close()
  defer db_info.Close()
  defer db_count.Close()
  defer db_breakpad_info.Close()

  // if ro and wo are not used again, be sure to Close them.

  max_num := 0
  s, _ := db_id.Get(ro, []byte("MAX_NUM"))
  if s != nil {
    //log.Println("MAX_NUM: ", string(s))
    num_val, err := strconv.Atoi(string(s))
    if err != nil {
      log.Println("MAX_NUM error:", ver)
      return
    }
    max_num = num_val
  }
  max_num++

  s, _ = db_id.Get(ro, []byte(address))
  if s == nil {
    max_num_val := strconv.Itoa(max_num)
    db_id.Put(wo, []byte("MAX_NUM"), []byte(max_num_val))
    db_id.Put(wo, []byte(address), []byte(max_num_val))
    db_id.Put(wo, []byte(max_num_val), []byte(address))
  }

  s, _ = db_info.Get(ro, []byte(address))
  if s == nil {
    log.Println("db_info Put:", address)
    db_info.Put(wo, []byte(address), []byte(info))
    db_breakpad_info.Put(wo, []byte(address), []byte(break_pad_info))
  }

  s, _ = db_count.Get(ro, []byte(address))
  if s == nil {
    log.Println("db_info Put:", address)
    db_count.Put(wo, []byte(address), []byte("1"))
  } else {
    num_val, err := strconv.Atoi(string(s))
    if err != nil {
      log.Println("db_count error:", ver)
      return
    }
    num_val++
    address_num_val := strconv.Itoa(num_val)
    db_count.Put(wo, []byte(address), []byte(address_num_val))
  }
}

func GetListInfoDB(pro string, ver string) string {
  opts := levigo.NewOptions()
  defer opts.Close()
  opts.SetCache(levigo.NewLRUCache(3 << 30))
  opts.SetCreateIfMissing(true)
  db_id, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_id.db", opts)
  db_info, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_info.db", opts)
  db_count, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_count.db", opts)
  db_breakpad_info, _ := levigo.Open("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", opts)
  defer db_id.Close()
  defer db_info.Close()
  defer db_count.Close()
  defer db_breakpad_info.Close()

  ro := levigo.NewReadOptions()
  defer ro.Close()

  return_val := "<html>\n<body>\n<table border=\"1\">\n"
  return_val = return_val + "<tr>\n"
  return_val = return_val + "<th align=\"left\">ID</th>\n"
  return_val = return_val + "<th align=\"right\">ADDRESS</th>\n"
  return_val = return_val + "<th align=\"right\">COUNT</th>\n"
  return_val = return_val + "<th align=\"center\">INFO</th>\n"
  return_val = return_val + "<th align=\"center\">BREAKPAD</th>\n"
  return_val = return_val + "</tr>\n"
  max_num_val, _ := db_id.Get(ro, []byte("MAX_NUM"))
  if max_num_val != nil {
    num_val, err := strconv.Atoi(string(max_num_val))
    if err != nil {
      log.Println("MAX_NUM error:", ver)
      return ""
    }
    max_num := num_val

    for i := 1; i <= max_num; i++ {
      id_val := strconv.Itoa(i)
      address, _ := db_id.Get(ro, []byte(id_val))

      if address != nil {

        info, _ := db_info.Get(ro, []byte(address))
        count, _ := db_count.Get(ro, []byte(address))
        info_pad, _ := db_breakpad_info.Get(ro, []byte(address))
        return_val = return_val + "<tr>\n"
        return_val = return_val + "<th align=\"left\">" + id_val + "</th>\n"
        return_val = return_val + "<th align=\"right\">" + string(address) + "</th>\n"
        return_val = return_val + "<th align=\"right\">" + string(count) + " </th>\n"
        return_val = return_val + "<th align=\"left\">" + string(info) + " </th>\n"
        return_val = return_val + "<th align=\"left\">" + string(info_pad) + " </th>\n"
        return_val = return_val + "</tr>\n"
      }
    }
  }
  return_val = return_val + "</table>\n</body>\n</html>"
  return return_val
}

func VerInfoDB(pro string) string {
  opts := levigo.NewOptions()
  defer opts.Close()
  opts.SetCache(levigo.NewLRUCache(3 << 30))
  opts.SetCreateIfMissing(true)
  db_version, _ := levigo.Open("./"+pro+"/version_id.db", opts)
  defer db_version.Close()

  ro := levigo.NewReadOptions()
  defer ro.Close()
  ro.SetFillCache(false)
  it := db_version.NewIterator(ro)
  defer it.Close()
  it.SeekToFirst()

  return_val := "<html>\n<body>\n"
  for it = it; it.Valid(); it.Next() {
    s := string(it.Key())
    return_val = return_val + "<a href=\"?par=get&ver=" + s + "\">" + s + "</a><br>"
  }
  return_val = return_val + "</body>\n</html>"
  return return_val
}
