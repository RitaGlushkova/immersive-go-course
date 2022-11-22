for i in {0..20}
 do 
    curl 'http://localhost:6060/debug/pprof/heap' -o "/Users/margaritaglushkova/pprof/pprof.batch-processing.alloc_objects.alloc_space.inuse_objects.inuse_space.007${i}.pb.gz" 
    sleep 0.5
done